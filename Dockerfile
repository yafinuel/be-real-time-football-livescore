FROM golang:1.27.3-alpine

WORKDIR /app

# Menyalin file dependency terlebih dahulu
COPY go.mod ./

# Download dependency jika sudah tersedia
RUN go mod download

# Source code akan di-mount melalui Docker Compose
CMD ["go", "run", "./cmd/api"]