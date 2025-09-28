# Stage 1: Build/Compile
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Opsional, tapi baik untuk caching layer
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build main app
RUN go build -o main ./cmd

# Build migration binary
RUN go build -o migrate ./cmd/migrate

# -----------------------------------------------------

# Stage 2: Production/Final Image
FROM alpine:3.22

WORKDIR /root/
RUN apk add --no-cache ca-certificates

# Salin semua aset yang dibutuhkan dari stage builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/cmd/migrate/migrations ./cmd/migrate/migrations/

EXPOSE 8080

# Perintah default saat kontainer dijalankan
CMD ["./main"]