# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o orders-api .

# Runtime stage
FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/orders-api .

RUN mkdir -p /app/data

EXPOSE 8002

ENV DB_PATH=/app/data/orders.db
ENV PORT=8002

CMD ["./orders-api"]
