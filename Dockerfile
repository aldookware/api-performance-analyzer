FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the analyzer CLI
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o analyzer ./cmd/analyzer

# Final stage
FROM alpine:latest

# Install git (needed for GitHub Actions context)
RUN apk --no-cache add ca-certificates git

WORKDIR /root/

# Copy the analyzer binary
COPY --from=builder /app/analyzer .

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
