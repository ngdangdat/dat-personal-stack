# Step 1: Build the Go binary
FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY services/ ./services/

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o node_agent cmd/node_agent/main.go

# Step 2: Minimal production container
FROM alpine:3.20.3

RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser

# Copy compiled binary from builder
COPY --from=builder /app/server .
COPY --from=builder /app/node_agent .

# Change ownership of binary files
RUN chown -R appuser:appgroup /home/appuser

USER appuser

EXPOSE 8080
EXPOSE 8081

CMD ["./server"]
