# -------- Stage 1: Build --------
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service

# -------- Stage 2: Runtime --------
FROM debian:bookworm-slim

WORKDIR /app

# Create non-root user
RUN useradd -m appuser

COPY --from=builder /app/auth-service .

USER appuser

EXPOSE 8000

CMD ["./auth-service"]
