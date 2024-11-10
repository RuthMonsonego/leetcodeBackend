# Dockerfile
FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

# עדכון תעודות SSL במיכל
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

# הגדרת GOPROXY למנוע בעיות של חיבור
ENV GOPROXY=https://proxy.golang.org,direct

# הורדת החבילות
RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

CMD ["/app/main"]
