FROM golang:1.24-alpine AS builder

RUN apk add --no-cache build-base git

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app cmd/api/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates && \
    adduser -D -u 1000 appuser

ENV APP_PORT=8080
WORKDIR /app

COPY --from=builder /go/src/app/app .
COPY --from=builder /go/src/app/api ./api

RUN chown -R appuser:appuser /app
USER appuser

EXPOSE ${APP_PORT}

HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${APP_PORT}/health || exit 1

CMD ["./app"]