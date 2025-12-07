FROM golang:1.24-alpine AS builder

WORKDIR /app

# Download deps first for better layer caching
COPY task-m-b/go.mod task-m-b/go.sum ./
RUN go mod download

# Copy the rest of the sources
COPY . .

# Build a static binary for linux
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/app ./cmd/app

FROM alpine:3.20

WORKDIR /root

# Copy the binary and container config
COPY --from=builder /app/app /root/app
COPY task-m-b/config/container.yaml /root/config/container.yaml

ENV CONFIG_PATH=/root/config/container.yaml

EXPOSE 80

CMD ["./app"]

