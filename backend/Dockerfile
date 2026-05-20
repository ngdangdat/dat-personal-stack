# Step 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod ./
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY services/ ./services/

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server cmd/server/main.go

# Step 2: Minimal production container
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy compiled binary from builder
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
