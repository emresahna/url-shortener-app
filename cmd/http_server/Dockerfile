FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY ../../go.mod go.sum ./

RUN go mod download

COPY ../.. .

RUN go build -o url-app ./cmd/http_server/main.go

FROM alpine:edge

COPY --from=builder /app/url-app .
COPY --from=builder /app/configs/ssl/ .

CMD ["./url-app"]