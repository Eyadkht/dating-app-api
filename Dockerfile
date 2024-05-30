# This Dockerfile is used to build the application using the golang:1.22.3-alpine base image.
# Using a builder image allows us to separate the build environment from the final image.
# This approach helps to keep the final image lightweight and helps in caching.
FROM golang:1.22.3-alpine as builder

# Set current directory
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 go build -o muzz-dating ./cmd/server

FROM alpine:latest

WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/muzz-dating .

# The TCP port the application is going to listen on by default.
EXPOSE 8888

# Run
CMD ["./muzz-dating"]