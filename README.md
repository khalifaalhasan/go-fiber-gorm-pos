# 🚀 Bangga Punya Web - Backend API

> Production-ready backend service built with **Clean Architecture**
> principles for a single-tenant business system.

------------------------------------------------------------------------

## 🎯 Development Objectives

Proyek ini dikembangkan dengan beberapa tujuan utama:

-   **Adopsi Clean Architecture**\
    Memisahkan logika bisnis (Service) dari akses data (Repository) dan
    pengiriman data (Controller).

-   **Single-Tenant Focus**\
    Mengoptimalkan performa dan keamanan data eksklusif untuk satu
    toko/entitas bisnis.

-   **Automatic Data Integrity**\
    Implementasi Automatic Slugging (menggunakan GORM Hooks) dan
    validasi data yang ketat.

-   **Developer Experience**\
    Struktur folder yang modular untuk memudahkan kolaborasi tim di masa
    depan (Visi Agency Bangga Punya Web).

------------------------------------------------------------------------

## 🛠️ Tech Stack

-   **Language:** Go (Golang) 1.2x\
-   **Web Framework:** Fiber (Express-like performance for Go)\
-   **ORM:** GORM (PostgreSQL)\
-   **Authentication:** JWT (JSON Web Token)\
-   **Validation:** Go-Playground Validator\
-   **Database:** PostgreSQL

------------------------------------------------------------------------

## 📂 Project Structure

    .
    ├── config/      # Database & Environment configuration
    ├── controller/  # Delivery layer (HTTP Request & Response)
    ├── middleware/  # JWT Protection & Security
    ├── model/       # Domain Entities & Data Contracts (DTO/Interface)
    ├── repository/  # Data Access Layer (GORM Queries)
    ├── routes/      # Modular Route Definitions
    ├── service/     # Business Logic Layer
    ├── utils/       # Helper functions (JWT, Logger, Validator)
    └── main.go      # Application Entry Point

------------------------------------------------------------------------

## 🚀 Key Features

-   ✅ Secure Authentication: Register & Login dengan enkripsi Bcrypt\
-   ✅ Category Management: CRUD kategori produk dengan fitur auto-slug\
-   ✅ Product Management: Manajemen menu lengkap dengan sistem promo
    dan harga normal\
-   ✅ Public API: Endpoint katalog menu khusus untuk pelanggan (SEO
    Friendly Slugs)\
-   ✅ Request Validation: Validasi input otomatis sebelum masuk ke
    database

------------------------------------------------------------------------

## 🏁 Quick Start

### 1️⃣ Clone Repository

``` bash
git clone <your-repository-url>
cd <project-folder>
```

### 2️⃣ Setup Environment

Pastikan file `.env` sudah terkonfigurasi:

    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=your_user
    DB_PASS=your_password
    DB_NAME=your_database
    JWT_SECRET=your_secret_key

### 3️⃣ Install Dependencies

``` bash
go mod tidy
```

### 4️⃣ Run Application

``` bash
go run main.go
```

Akses API di:

    http://localhost:8080

------------------------------------------------------------------------

## ❤️ Author

Dibuat dengan ❤️ di Palembang oleh **Khalifa Al Hasan** 🚀☕

------------------------------------------------------------------------

## 📌 Vision

Backend ini merupakan bagian dari visi besar **Bangga Punya Web
Agency**\
untuk membangun sistem yang scalable, secure, dan siap production.
