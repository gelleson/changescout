# Use an official Go runtime as a parent image
FROM golang:1.23-alpine AS builder

# Set the working directory in the container
WORKDIR /app

# Copy the go.mod and go.sum files into the container
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the application code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# Use a slimmer image for the final runtime environment
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy only the built application binary from the builder stage
COPY --from=builder /app/main ./

# Expose the port
EXPOSE 3311

# Define the entrypoint
ENTRYPOINT ["./main"]
