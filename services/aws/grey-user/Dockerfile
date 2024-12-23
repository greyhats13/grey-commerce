# Path: grey-user/Dockerfile

# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install necessary packages
RUN apk update && apk add --no-cache git

# Copy go.mod and go.sum for dependency resolution
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o grey-user ./cmd/main.go

# Final stage
FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/grey-user /app/grey-user

# Service will listen on PORT (set via env or config)
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/grey-user"]
