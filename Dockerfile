FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download && go build -o urlshortener .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/urlshortener .

EXPOSE 8080

CMD ["./urlshortener"]
