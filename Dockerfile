# docker build -t go-inventory-v2 .
# docker run -dit --name go-inventory-v2 -p 7100:7100 go-inventory-v2

FROM golang:1.25 AS builder

RUN apt-get update && apt-get install bash && apt-get install -y --no-install-recommends ca-certificates

WORKDIR /app
COPY . .
RUN go mod tidy

WORKDIR /app
RUN go build -o go-inventory-v2 -ldflags '-linkmode external -w -extldflags "-static"'

FROM alpine

WORKDIR /app
COPY --from=builder /app/go-inventory-v2 .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app/go-inventory-v2"]