# Stage 1: Builder
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install root CA certificates for HTTPS routing
RUN apk add --no-cache ca-certificates

COPY . .
RUN go mod init panovista-core && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o panovista .

# Stage 2: The Zero-Surface Vault
FROM scratch

# Import the certificates from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/panovista /panovista
ENTRYPOINT ["/panovista"]