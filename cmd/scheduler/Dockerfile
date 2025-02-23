FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY ../../go.mod go.sum ./

RUN go mod download

COPY ../.. .

RUN go build -o scheduler ./cmd/scheduler/main.go

FROM alpine:edge

COPY --from=builder /app/scheduler .

EXPOSE 3000

CMD ["./scheduler"]