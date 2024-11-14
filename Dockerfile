FROM golang:1.23-alpine

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Install bash
RUN apk add --no-cache bash

# Copy the application files
COPY . .

# Build the application
RUN go build -o /main ./cmd/main.go

# Expose the application port
EXPOSE 8080

# Run the built application
CMD ["/cmd/main"]
