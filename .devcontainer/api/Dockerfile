FROM mcr.microsoft.com/devcontainers/go:1.22-bookworm

# Install Air for Go hot-reload
RUN go install github.com/air-verse/air@latest

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.57.2

# Install mysql client
RUN apt-get update && apt-get install -y default-mysql-client \
    && apt-get clean && rm -rf /var/lib/apt/lists/*
