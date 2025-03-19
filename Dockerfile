# Stage 1: Build
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main cmd/main.go

# Stage 2: Run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env ./.env

# Add a health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
