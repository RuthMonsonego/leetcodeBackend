FROM golang:1.20

WORKDIR /app

# Copy user Go code into the container
COPY user_code.go .

# Build the Go application
RUN go build -o main user_code.go

# Run the Go application
CMD ["./main"]
