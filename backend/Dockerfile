FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o fintech-wallet .
RUN go build -o worker ./cmd/worker
RUN go build -o outbox ./cmd/outbox


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/fintech-wallet .
COPY --from=builder /app/worker .
COPY --from=builder /app/outbox .

EXPOSE 8081

CMD ["./fintech-wallet"]