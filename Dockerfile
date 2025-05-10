# Stage 1: Build
FROM golang:1.23 AS builder

# Making working dir
WORKDIR /app

# Installing linter
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.58.0

# Copy go.mod and go.sum and installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy code
COPY . .

# Building binary
RUN go build -o main .

# Stage 2: Run
FROM gcr.io/distroless/base-debian12

# Installing dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Making working dir
WORKDIR /app

# Copy binary file from first stage
COPY --from=builder /app/main .

# App runner
CMD ["./main"]
