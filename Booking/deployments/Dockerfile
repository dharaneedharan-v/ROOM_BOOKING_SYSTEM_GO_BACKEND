
# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Compiles your entry point correctly from the root context
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/...

# Run stage
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache tzdata

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs


RUN mkdir logs
CMD ["./main"]
