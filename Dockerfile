FROM golang:1.21-alpine

RUN apk update && \
    apk add --no-cache \
    docker-cli \
    curl \
    git \
    ca-certificates \
    && update-ca-certificates

RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/

WORKDIR /app

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/main.go

EXPOSE 8080
CMD ["./main"]