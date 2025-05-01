FROM golang:1.24-alpine AS builder
RUN apk add --no-cache gcc musl-dev git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o expense-tracker

FROM alpine:latest
RUN apk add --no-cache ca-certificates sqlite tzdata

# Create app structure with correct permissions
RUN mkdir -p /app/data && \
    chown -R nobody:nobody /app && \
    chmod -R 775 /app

WORKDIR /app
USER nobody

COPY --from=builder /app/expense-tracker .
COPY --from=builder /app/views ./views
COPY --from=builder /app/static ./static

ENV DB_PATH=/app/data/transactions.db \
    TZ=UTC

EXPOSE 3000
CMD ["./expense-tracker"]