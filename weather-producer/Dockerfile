# Use the official Golang image from the Docker Hub
FROM golang:1.20 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 go build -o weather-producer .

# Start a new stage from scratch
FROM debian:bullseye-slim

# Set the Current Working Directory inside the container
WORKDIR /root/

# Install CA certificates
RUN apt-get update && apt-get install -y ca-certificates

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/weather-producer .

# Expose port 8080 to the outside world (if needed)
EXPOSE 8080

# Command to run the executable
CMD ["./weather-producer"]
