FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o url ./cmd/main.go

FROM alpine:edge

COPY --from=builder /app/url .
COPY --from=builder /app/configs/ssl/ .

CMD ["./url"]