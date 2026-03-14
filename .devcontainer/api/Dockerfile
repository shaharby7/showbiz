FROM mcr.microsoft.com/devcontainers/go:1.22-bookworm

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.57.2

# Install mysql client (remove stale yarn repo from base image first)
RUN rm -f /etc/apt/sources.list.d/yarn.list \
    && apt-get update && apt-get install -y default-mysql-client \
    && apt-get clean && rm -rf /var/lib/apt/lists/*
