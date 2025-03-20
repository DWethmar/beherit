FROM mcr.microsoft.com/devcontainers/go:dev-1.24

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.64.6

# install cobra cli
RUN go install github.com/spf13/cobra-cli@latest

# allow user to install go tools
RUN sudo chown -R vscode:vscode /go/pkg

# Ensure githooks directory exists
RUN mkdir -p /workspace/.githooks

# Set Git hooks path
RUN git config --global core.hooksPath /workspace/.githooks

CMD ["sleep", "infinity"]