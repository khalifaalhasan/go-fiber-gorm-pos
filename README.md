# 🚀 Bangga Punya Web - POS Backend API

> **Enterprise-grade POS Backend** built with Golang Fiber, following Clean Architecture principles and high-performance patterns.

---

## 🎯 Technical Excellence & Features

Proyek ini tidak hanya sekadar CRUD, melainkan mengimplementasikan standar industri untuk reliabilitas dan skalabilitas:

### 1. 🛡️ Idempotency Handling
Menjamin bahwa request yang sama (misal: double-click checkout) tidak akan memproses order dua kali.
- **Distributed Locking**: Menggunakan Redis `SETNX` untuk mengunci proses yang sedang berjalan.
- **Permanent Caching**: Hasil response sukses disimpan di Database, sehingga jika client mengirim ulang dengan `Idempotency-Key` yang sama, server akan mengembalikan hasil yang sama tanpa membebani logika bisnis.

### 2. ⚡ N+1 Query Resolution
Menghindari masalah performa database klasik di mana 1 request menghasilkan puluhan query tambahan.
- **Eager Loading**: Menggunakan GORM `.Preload()` untuk mengambil relasi (seperti Order Items & Products) dalam satu batch query yang efisien.

### 3. 🛡️ Safe Concurrency & Deadlock Prevention
Sistem ini menangani ribuan transaksi bersamaan dengan aman:
- **Pessimistic Locking**: Menggunakan `SELECT ... FOR UPDATE` untuk menjamin konsistensi stok dan nomor antrean.
- **Deadlock Avoidance**: Mengimplementasikan **Global Sorting Strategy** pada ID produk sebelum akuisisi lock untuk mencegah circular wait.

### 4. 🐳 Docker Multistage Build
Optimasi container untuk deployment production:
- **Build Stage**: Menggunakan image Golang Alpine untuk kompilasi.
- **Runtime Stage**: Hanya menyertakan binary statis di atas image Alpine murni (~5MB), menghasilkan image yang sangat ringan, aman, dan cepat di-deploy.

---

## 🛠️ Tech Stack

- **Framework**: [Fiber v2](https://gofiber.io/) (Fastest Go Web Framework)
- **DB & Cache**: PostgreSQL & Redis
- **ORM**: GORM
- **Docs**: Swagger (Swaggo)
- **Validation**: Go-Validator v10

---

## 📖 API Documentation & Testing

### 🟢 Swagger UI
Dokumentasi interaktif dapat diakses saat aplikasi berjalan:
1. Jalankan aplikasi: `go run cmd/api/main.go` atau `docker-compose up`
2. Buka: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### 🟠 Postman Collection
Untuk pengujian flow End-to-End yang lebih komprehensif:
1. Import file `Postman_Collection_EndToEnd.json` ke Postman.
2. Pastikan environment variable `base_url` mengarah ke `http://localhost:8080`.
3. Gunakan folder `1. Authentication` untuk mendapatkan token, yang akan otomatis tersimpan untuk request berikutnya.

---

## 🏁 Quick Start

### 1️⃣ Clone & Run with Docker (Recommended)
```bash
docker-compose up -d --build
```

### 2️⃣ Manual Setup
1. Copy `.env` dan sesuaikan kredensial database.
2. `go mod tidy`
3. Generate docs: `swag init -g cmd/api/main.go`
4. Jalankan: `go run cmd/api/main.go`

---

## ❤️ Author
Dibuat dengan ❤️ oleh **Khalifa Al Hasan** 🚀☕\
Part of **Bangga Punya Web Agency** Vision.
