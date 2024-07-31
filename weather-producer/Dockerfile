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

# Load environment variables from the .env file into system-wide environment variables
RUN export $(grep -v '^#' /root/.env | xargs) && \
    echo 'ENV_FILE=/root/.env' >> /root/.profile && \
    echo 'for var in $(grep -v '^#' /root/.env | awk -F= '{print $1}'); do export $var=$(grep ^$var= /root/.env | cut -d= -f2-); done' >> /root/.profile

# Expose port 8080 to the outside world (if needed)
EXPOSE 8080

# Command to run the executable
CMD ["./weather-producer"]
