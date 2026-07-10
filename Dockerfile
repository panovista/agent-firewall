# ---------------------------------------------------
# STAGE 1: COMPILATION
# ---------------------------------------------------
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY main.go .

# Initialize a temporary Go module and compile as a statically linked binary
RUN go mod init panovista-core && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o panovista .

# ---------------------------------------------------
# STAGE 2: IMMUTABLE SCRATCH CONTAINER
# ---------------------------------------------------
FROM scratch

# Copy only the compiled machine-code binary from Stage 1
COPY --from=builder /app/panovista /panovista

# Force the container to execute the binary directly
ENTRYPOINT ["/panovista"]