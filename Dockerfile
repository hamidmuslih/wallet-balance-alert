# Stage 1: Build the Go application
FROM golang:1.22.4-alpine@sha256:ace6cc3fe58d0c7b12303c57afe6d6724851152df55e08057b43990b927ad5e8 AS builder

WORKDIR /app

# Copy the application code and download dependencies
COPY . .
RUN go mod download

# Build the Go application
RUN go build -o wallet-checker

# Stage 2: Create a minimal image with the built binary
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/wallet-checker .

# Run the binary
ENTRYPOINT ["./wallet-checker"]
