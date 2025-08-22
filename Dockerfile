FROM golang:1.24-alpine AS builder

RUN go version

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY .env .
COPY config/local.yaml ./config/

COPY . .
RUN go build -o main ./cmd/app

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./cmd/app

FROM alpine:3.20

WORKDIR /root/

COPY --from=0 /app/app .

CMD ["./app"]

