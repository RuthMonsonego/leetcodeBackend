# Dockerfile for the Go application
FROM golang:1.19-alpine

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./

# Copy all files
COPY cmd/main.go ./
COPY config/config.go ./
COPY controllers/question_controller.go ./
COPY models/question.go ./
COPY repositories/question_repository.go ./

RUN apk add --no-cache ca-certificates

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o /main

# Expose the application port
EXPOSE 8080

# Run the built application
CMD ["/main"]
