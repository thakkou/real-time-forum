# Build stage

FROM golang:1.24-alpine AS builder

LABEL authors="erezzoug thakkou halhyane herraba" version="1.0"
LABEL authors="thakkou herraba" version="2.0"

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o forum .

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/forum .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/database/schema.sql ./database/schema.sql

EXPOSE 8080

ENTRYPOINT ["./forum"]