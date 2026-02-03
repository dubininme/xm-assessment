FROM golang:1.24-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@v1.62.0

RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1

COPY go.mod go.sum ./
RUN go mod download

EXPOSE 8080

# Simple go run - restart manually with docker-compose restart app
CMD ["go", "run", "./cmd/api/main.go"]
