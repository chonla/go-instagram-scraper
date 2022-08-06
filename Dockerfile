FROM golang:1.19.0-alpine3.16 AS builder

RUN apk add build-base

WORKDIR /app
COPY . .
RUN go mod download && \
    go build -o instagram-scraper

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/instagram-scraper /app/instagram-scraper
CMD ["/app/instagram-scraper"]