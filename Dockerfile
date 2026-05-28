# Stage 1: Build the binary inside a Go environment
FROM golang:1.26-alpine AS builder
WORKDIR /app

# Copy dependency manifests and download them
COPY go.mod go.sum ./
RUN go mod download

# Copy the source files and build a statically linked native binary
COPY main.go hello_controller.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app-server .

# Stage 2: Create the minimal runtime container
FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/app-server .

# Expose the app port matching your main.go setup
EXPOSE 8081

# Run the application
CMD ["./app-server"]