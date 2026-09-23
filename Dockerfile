# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o image-processing-service .


# Runtime stage
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/image-processing-service .

EXPOSE 8000

CMD ["./image-processing-service"]
