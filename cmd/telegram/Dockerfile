FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o telegram ./cmd/telegram

# Этап 2: Финальный образ
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Копируем бинарник из этапа сборки
COPY --from=builder /app/telegram .

# Запускаем
CMD ["./telegram"]