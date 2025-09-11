# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env .
COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./server"]