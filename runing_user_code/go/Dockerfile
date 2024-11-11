# Use the official Golang image, version 1.20
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy user Go code into the container
COPY user_code.go .

# Build the Go application
RUN go build -o main user_code.go

# Run the Go application with user input
CMD ["./main"]
