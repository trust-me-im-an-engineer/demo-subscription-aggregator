FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /subscription-aggregator ./cmd/server

FROM alpine:latest

WORKDIR /app
COPY --from=builder /subscription-aggregator .
COPY migrations ./migrations

EXPOSE 8080

CMD ["./subscription-aggregator"]