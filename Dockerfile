# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubectl-pilot .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates kubectl

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/kubectl-pilot .

# Copy example config
COPY examples/config.yaml .k8s-pilot.yaml

# Set environment variables
ENV K8S_PILOT_AI_PROVIDER=mock

ENTRYPOINT ["./kubectl-pilot"]
