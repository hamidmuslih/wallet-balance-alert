# Stage 1: Build the Go application
FROM golang:1.22.4-alpine@sha256:ace6cc3fe58d0c7b12303c57afe6d6724851152df55e08057b43990b927ad5e8 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o wallet-checker

# Stage 2: Create a minimal image with the built binary
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/wallet-checker .

# Run the binary
ENTRYPOINT ["./wallet-checker"]
