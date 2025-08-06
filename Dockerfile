# Build Stage
FROM golang:1.24.5 AS builder

WORKDIR /app

# Copy Go source code
COPY . .

# Build the Go binary
RUN go build -o webpage-analyzer-go main.go

# Final Stage
FROM debian:bullseye-slim

WORKDIR /app

# Copy binary and required files from builder
COPY --from=builder /app/webpage-analyzer .
COPY --from=builder /app/app.yaml .

# Expose port (adjust if needed)
EXPOSE 8080

# Run the binary
CMD ["./webpage-analyzer"]
