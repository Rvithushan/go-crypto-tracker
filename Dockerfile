# Stage 1: Build base using Goland Alpine image 

# Stage 1: Build stage using Golang Alpine image
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy Go module file
COPY go.mod ./

# Copy all source code (main.go, index.html)
COPY . .

# Compile the Go application into a standalone static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o crypto-tracker .

# Stage 2: Minimal runner stage using Alpine Linux
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates so Go can make HTTPS requests to CoinGecko API
RUN apk --no-cache add ca-certificates

# Copy compiled binary and HTML template from builder stage
COPY --from=builder /app/crypto-tracker .
COPY --from=builder /app/index.html .

EXPOSE 8080

CMD ["./crypto-tracker"]
