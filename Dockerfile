# Step 1: Base image for building the Go application
FROM golang:1.23.1-alpine AS builder

# Install necessary tools (git, curl)
RUN apk add --no-cache git curl

# Install air (a hot-reload tool for Go)
RUN go install github.com/air-verse/air@latest

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application code
COPY . .

# Build the Go application
RUN go build -o server .

# Set the default command to run air for hot-reloading
CMD ["air"]
