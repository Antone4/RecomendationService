# Stage 1: Build
FROM golang:1.24.2-alpine AS builder

# Устанавливаем зависимости, нужные для сборки Go-проектов
RUN apk add --no-cache git libc6-compat

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o recommendation .

# Stage 2: Run
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/recommendation .

EXPOSE 8080

CMD ["./recommendation"]
