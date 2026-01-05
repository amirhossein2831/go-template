# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# install proto buf generator
RUN apk update && \
    apk add --no-cache \
    protobuf \
    protobuf-dev \
    build-base
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy dependency files
ENV GOPROXY=https://goproxy.io,direct
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY ./ ./

# Generate the Go code from .proto files by running your script.
RUN chmod +x ./scripts/proto-gen.sh && ./scripts/proto-gen.sh

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/app ./cmd/main.go

# Stage 2: Runtime
FROM alpine:3.19
RUN apk add --no-cache ca-certificates

# Create binary directory
WORKDIR /app

# Copy all binaries from builder
COPY --from=builder /app/bin/app /app/app
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/internal/database/migrations /app/internal/database/migrations

ENTRYPOINT ["/app/app"]