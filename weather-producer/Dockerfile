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

# Copy the .env file into the container
COPY .env /root/.env

# Create a script to load environment variables from the .env file
RUN echo '#!/bin/sh\n\
set -a\n\
. /root/.env\n\
set +a\n\
exec "$@"' > /root/entrypoint.sh && chmod +x /root/entrypoint.sh

# Expose port 8080 to the outside world (if needed)
EXPOSE 8080

# Use the entrypoint script to load environment variables and run the executable
ENTRYPOINT ["/root/entrypoint.sh"]
CMD ["./weather-producer"]
