# ==========================================
# Build Stage
# ==========================================
FROM golang:1.22-alpine AS builder

# Set security environment flags and prevent CGO dependencies
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /workspace

# Cache Go modules layer
COPY go.mod ./
# RUN go mod download # (if external modules exist)

# Copy application source code
COPY src/ ./src/

# Build statically-linked, stripped binary
ARG VERSION=1.0.0
RUN go build -ldflags="-s -w -X main.version=${VERSION}" -o /workspace/bin/core-service ./src

# ==========================================
# Production Distroless Secure Stage
# ==========================================
# Google Distroless static non-root (Zero CVEs, no package manager, no shell)
FROM gcr.io/distroless/static-debian12:nonroot

LABEL maintainer="DevSecOps Platform Team" \
      project="multi-cloud-devsecops-platform" \
      security.policy="zero-trust"

WORKDIR /

# Copy executable from builder
COPY --from=builder /workspace/bin/core-service /core-service

# Explicit non-root user (nonroot:nonroot UID:GID 65532:65532)
USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/core-service"]
