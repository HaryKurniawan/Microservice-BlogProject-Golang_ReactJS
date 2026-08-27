# Stage 1: Build binary Go
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy semua source code
COPY . .

# Download dependencies dan build binary
RUN go mod tidy && go build -o api-gateway .

# Stage 2: Runtime image yang ringan
FROM alpine:latest

WORKDIR /app

# Copy binary hasil build dari stage sebelumnya
COPY --from=builder /app/api-gateway .

EXPOSE 8080

CMD ["./api-gateway"]
