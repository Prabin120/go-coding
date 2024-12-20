# Step 1: Base image for building the Go application
FROM golang:1.23.1-alpine AS builder

# Install necessary tools (git, curl, and compilers/interpreters for various languages)
RUN apk add --no-cache git curl \
    python3 py3-pip \
    gcc g++ musl-dev \
    openjdk17 \
    nodejs npm

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application code
COPY . .

# Build the Go application
RUN go build -o server main.go

# Set the default command to run the application
CMD ["./server"]