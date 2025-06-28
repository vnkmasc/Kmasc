
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server


FROM alpine:latest

# Cài chứng chỉ CA cho HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary từ stage builder sang stage chạy
COPY --from=builder /app/app .

# Copy các tài nguyên cần thiết nếu bạn sử dụng
COPY --from=builder /app/internal/config ./internal/config
COPY --from=builder /app/wallet ./wallet
COPY --from=builder /app/keystore ./keystore
COPY --from=builder /app/.env .env

# Cổng mà app listen
EXPOSE 8080

# Chạy app
CMD ["./app"]
