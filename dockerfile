FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o tradeverse ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/tradeverse .

EXPOSE 8080

CMD ["./tradeverse"]
