# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY go.mod ./ 

RUN go mod download

COPY main.go .
# Build a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notify-app .

# Stage 2: Runtime
FROM alpine:latest
# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
# Copy only the compiled binary from the builder stage
COPY --from=builder /app/notify-app .
# Switch to the non-root user
USER appuser

EXPOSE 8080
CMD ["./notify-app"]
