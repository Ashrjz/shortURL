# Stage 1: Build
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY api/ ./api/

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./api

# Stage 2: Run (small final image)
FROM alpine:latest

WORKDIR /app

# Copy only the binary from builder stage
COPY --from=builder /app/main .

ENV GIN_MODE=release

EXPOSE 8080

CMD ["./main"]