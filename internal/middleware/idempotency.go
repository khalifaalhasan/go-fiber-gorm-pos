package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"

	"go-fiber-pos/internal/core"
	"go-fiber-pos/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// IdempotencyMiddleware memastikan bahwa request dengan Idempotency-Key yang sama
// tidak diproses berulang kali secara concurrent (membantu mencegah double booking).
func IdempotencyMiddleware(redisClient *redis.Client, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil Idempotency-Key dari Header
		idempotencyKey := c.Get("Idempotency-Key")
		if idempotencyKey == "" {
			// Fallback ke X-Request-ID jika diinginkan (opsional)
			idempotencyKey = c.Get("X-Request-ID")
		}

		// Jika tidak ada key, abaikan middleware ini dan jalankan seperti biasa
		if idempotencyKey == "" {
			return c.Next()
		}

		ctx := context.Background()
		redisKey := "idempotency:lock:" + idempotencyKey

		// ==========================================
		// SCENARIO B: Cek Database (Permanent Record)
		// ==========================================
		// Jika request sebelumnya sudah selesai diproses, maka datanya ada di database.
		var record core.IdempotencyRecord
		err := db.First(&record, "key = ?", idempotencyKey).Error
		if err == nil {
			// Request sudah pernah berhasil dikerjakan sebelumnya.
			// Return respons yang persis sama dari DB (tanpa memproses order lagi).
			logger.Log.Infof("Idempotency hit from DB for key: %s", idempotencyKey)
			c.Set("X-Idempotency-Hit", "true")
			c.Status(record.StatusCode)
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return c.SendString(record.ResponseBody)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Errorf("Error checking idempotency DB: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error during idempotency check",
			})
		}

		// ==========================================
		// SCENARIO C/A: Try to acquire lock di Redis
		// ==========================================
		// SETNX (Set if Not Exists)
		// TTL pendek (30 detik) mengasumsikan proses order tidak lebih lama dari ini
		acquired, err := redisClient.SetNX(ctx, redisKey, "processing", 30*time.Second).Result()
		if err != nil {
			logger.Log.Errorf("Error setting idempotency lock in Redis: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error during distributed locking",
			})
		}

		// SCENARIO A: Jika acquired == false
		// Artinya request lain dengan Key yang sama SEDANG berjalan di node/thread lain
		if !acquired {
			logger.Log.Warnf("Idempotency conflict for key: %s. Request is currently in-flight.", idempotencyKey)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":          "Request is currently being processed. Please wait or try again later.",
				"idempotency_key": idempotencyKey,
			})
		}

		// SCENARIO C: Berhasil dapat lock. Kita lepaskan lock di akhir proses bagaimanapun caranya.
		defer func() {
			// Release lock di Redis
			if err := redisClient.Del(ctx, redisKey).Err(); err != nil {
				logger.Log.Errorf("Gagal menghapus kunci lisensi Redis %s: %v", redisKey, err)
			}
		}()

		// ==========================================
		// JALANKAN BUSINESS LOGIC
		// ==========================================
		err = c.Next()

		// ==========================================
		// SIMPAN HASIL KE DATABASE
		// ==========================================
		// Hanya simpan jika response adalah 2xx (Success)
		statusCode := c.Response().StatusCode()
		if statusCode >= 200 && statusCode < 300 {
			// Kita perlu membaca body dari fiber Response
			respBody := c.Response().Body()

			// Karena framework mengganti spasi saat JSON marshalling jika tidak hati2, pastikan nilainya valid
			compactBody := new(bytes.Buffer)
			if errJson := json.Compact(compactBody, respBody); errJson == nil {
				respBody = compactBody.Bytes()
			}

			newRecord := &core.IdempotencyRecord{
				Key:          idempotencyKey,
				ResponseBody: string(respBody),
				StatusCode:   statusCode,
			}

			// Simpan ke DB
			if dbErr := db.Create(newRecord).Error; dbErr != nil {
				logger.Log.Errorf("Gagal menyimpan record idempotency ke DB: %v", dbErr)
				// Jangan kembalikan error ke client karena core logic sudah sukses, cukup di log
			}
		}

		return err
	}
}
