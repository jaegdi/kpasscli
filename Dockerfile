# Use Red Hat UBI 7 as the base image
FROM registry.access.redhat.com/ubi7/ubi:latest

USER 0

# Install dependencies
RUN yum install -y \
    wget \
    tar \
    gcc \
    make \
    libX11-devel \
    glibc-devel \
    && yum clean all

# Install Go manually
ENV GO_VERSION=1.23.7
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# Set the working directory
WORKDIR /apprking directory
WORKDIR /app

# Copy the source code into the container
COPY . .
RUN rm go.mod go.sum \
    && go mod init kpasscli

# Build the kpasscli binary with updated C99 and GCC flags
ENV CGO_CFLAGS="-std=c99 -Wno-implicit-function-declaration" \
    CGO_ENABLED=0
RUN go mod tidy \
 && go build -v -o dist

# Copy the built binary to the local directory
CMD ["cp", "./kpasscli-ubi7", "./dist/kpasscli-ubi7"]

# Set the entrypoint to the built binary
ENTRYPOINT ["./kpasscli"]
