FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY ../../go.mod go.sum ./

RUN go mod download

COPY ../.. .

RUN go build -o api ./cmd/api/main.go

FROM alpine:edge

COPY --from=builder /app/api .
COPY --from=builder /app/configs/ssl/ .

CMD ["./api"]