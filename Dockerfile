# Build Stage
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go-api ./cmd/app/main.go

# Run Stage
FROM alpine:latest

# Add CA certificates for HTTPS support (if needed)
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /go-api /root/go-api

# Ensure the binary is executable
RUN chmod +x /root/go-api

# Expose the port the application will run on
EXPOSE 8081

# Set the entrypoint to the binary
CMD ["/root/go-api"]