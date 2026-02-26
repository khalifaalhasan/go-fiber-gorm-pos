package config

import (
	"context"
	"os"

	"go-fiber-pos/pkg/logger"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// ConnectRedis melakukan inisialisasi koneksi ke Redis Server
func ConnectRedis() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0, // Use default DB
	})

	// PING untuk memastikan test koneksi berhasil
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logger.Log.Fatalf("Gagal terhubung ke Redis: %v", err)
	}

	RedisClient = client
	logger.Log.Info("Berhasil terhubung ke Redis!")
}
