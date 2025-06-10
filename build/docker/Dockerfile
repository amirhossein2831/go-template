# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY ./ ./

# Build the application
# -ldflags="-s -w" reduces binary size by stripping debug info
# CGO_ENABLED=0 creates a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/ ./cmd/...

# Stage 2: Runtime
FROM alpine:3.19
RUN apk add --no-cache ca-certificates

# Create binary directory
WORKDIR /app

# Copy all binaries from builder
COPY --from=builder --chown=appuser:appgroup /app/bin/ /app/bin/
COPY --from=builder --chown=appuser:appgroup /app/configs/ /app/configs/
COPY --from=builder --chown=appuser:appgroup /app/script/entrypoint.sh /app/

# add permissions
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]