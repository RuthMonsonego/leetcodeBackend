FROM golang:1.23-alpine

WORKDIR /app

COPY wait-for-it.sh /wait-for-it.sh

# Copy go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

RUN apk add --no-cache bash

# Copy the application files
COPY . .

# Build the application
RUN go build -o /main ./cmd/main.go

# Expose the application port
EXPOSE 8080

# Run the built application
CMD ["/main"]
