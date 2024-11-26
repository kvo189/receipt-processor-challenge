# Use the official Go image as the base image
FROM golang:1.23.3

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod tidy

# Copy the rest of the application code into the container
COPY . .

# Build the Go application
RUN go build -o receipt-processor ./cmd/main.go

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./receipt-processor"]
