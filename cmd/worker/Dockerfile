FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY ../../go.mod go.sum ./

RUN go mod download

COPY ../.. .

RUN go build -o worker ./cmd/worker/main.go

FROM alpine:edge

COPY --from=builder /app/worker .

EXPOSE 3000

CMD ["./worker"]