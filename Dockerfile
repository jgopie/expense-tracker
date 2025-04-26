FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o expense-tracker

# Run the application
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/expense-tracker .

EXPOSE 3000
CMD ["./expense-tracker"]
