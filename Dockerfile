# Use the official Golang image
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod ./

# Download and install Go modules
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal base image for the final container
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the application runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
