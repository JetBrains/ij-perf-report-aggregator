FROM mcr.microsoft.com/devcontainers/base:ubuntu

# Install Node.js
ARG NODE_VERSION="20"
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash - \
    && apt-get update \
    && apt-get install -y nodejs

# Install Go
RUN apt-get update && apt-get install -y golang
ENV PATH=$PATH:/usr/local/go/bin:/go/bin
ENV GOPATH=/go

# Install common development tools
RUN apt-get update && apt-get install -y \
    git \
    curl \
    wget \
    zip \
    unzip \
    build-essential \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create go directory
RUN mkdir -p /go/bin /go/src /go/pkg
RUN chmod -R 777 /go

# Create a new user for the project
ARG USERNAME=ijperf

# Create user without sudo access
RUN useradd -m -s /bin/bash $USERNAME

# Install Vue CLI and pnpm globally as root
RUN npm install -g @vue/cli pnpm

# Set the default user
USER $USERNAME

# Set the working directory
WORKDIR /workspace
