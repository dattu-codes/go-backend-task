# Stage 1: Build stage using official Go alpine image
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy dependency manifests
COPY go.mod go.sum ./

# Download dependencies (cached layer)
RUN go mod download

# Copy all source files
COPY . .

# Build a statically linked executable
# CGO_ENABLED=0 makes the binary independent of C library versions
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Stage 2: Final release runtime using a minimal Alpine container
FROM alpine:latest  

# Install certificates for secure outbound network calls
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the compiled binary from the builder environment
COPY --from=builder /app/main .

# Expose server port
EXPOSE 3000

# Start application
CMD ["./main"]
