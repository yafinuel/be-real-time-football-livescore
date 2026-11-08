# ─────────────────────────────────────────────
# Stage 1: Development & Dependencies
# ─────────────────────────────────────────────
FROM golang:1.27.3-alpine AS dev

WORKDIR /app

# Salin file dependency
COPY go.mod go.sum ./
RUN go mod download

# Salin source code
COPY . .

EXPOSE 8080

CMD ["go", "run", "./cmd/api"]

# ─────────────────────────────────────────────
# Stage 2: Builder untuk Production
# ─────────────────────────────────────────────
FROM dev AS builder

# Compile binary statically (tanpa CGO, strip debug symbols agar ukuran kecil)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/api

# ─────────────────────────────────────────────
# Stage 3: Minimal Production Runtime (~15-20 MB)
# ─────────────────────────────────────────────
FROM alpine:3.20 AS prod

WORKDIR /app

# Install root CA certificates dan data timezone
RUN apk --no-cache add ca-certificates tzdata

# Salin binary hasil kompilasi
COPY --from=builder /app/server /app/server

EXPOSE 8080

CMD ["/app/server"]
