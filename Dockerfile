# ==========================================
# STAGE 1: Builder (Membangun Aplikasi)
# ==========================================
# Menggunakan image Golang resmi versi Alpine agar ringan
FROM golang:1.22-alpine AS builder

# Set working directory di dalam container
WORKDIR /app

# BEST PRACTICE: Copy file module duluan untuk caching layer.
# Jika tidak ada perubahan dependensi, Docker tidak akan download ulang.
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi Golang
# CGO_ENABLED=0 membuat binary statis murni (tidak bergantung pada library C OS host)
# Path disesuaikan dengan struktur folder Anda: cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-binary ./cmd/api/main.go

# ==========================================
# STAGE 2: Runner (Environment Production yang Bersih)
# ==========================================
# Menggunakan Alpine murni yang sangat kecil (~5MB)
FROM alpine:latest

WORKDIR /app

# Menginstal sertifikat SSL (penting jika API Anda hit endpoint eksternal/Midtrans)
RUN apk --no-cache add ca-certificates tzdata

# Copy HANYA file binary hasil build dari Stage 1
COPY --from=builder /app/api-binary .

# Expose port yang digunakan oleh Go Fiber (misal: 3000)
EXPOSE 3000

# Eksekusi aplikasi
CMD ["./api-binary"]