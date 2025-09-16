# Dynamic Builds
ARG BUILDER_IMAGE=golang:1.25-bookworm
ARG FINAL_IMAGE=debian:bookworm-slim

# Build stage
FROM --platform=${BUILDPLATFORM} ${BUILDER_IMAGE} AS builder

# Build Args
ARG GIT_REVISION=""

# Platform args
ARG TARGETOS
ARG TARGETARCH

# Ensure ca-certificates are up to date
RUN update-ca-certificates

# Use modules for dependencies
WORKDIR $GOPATH/src/go.rtnl.ai/honu

COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

# Copy package
COPY . .

# Build binary
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-X 'go.rtnl.ai/honu/pkg.GitVersion=${GIT_REVISION}'" \
    -o /go/bin/honudb \
    ./cmd/honudb

# Final Stage
FROM --platform=${BUILDPLATFORM} ${FINAL_IMAGE} AS final

LABEL maintainer="Rotational Labs <support@rotational.io>"
LABEL description="Honu distributed database replica instance"

# Ensure ca-certificates are up to date
RUN set -x && apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage
COPY --from=builder /go/bin/honudb /usr/local/bin/honudb

CMD [ "/usr/local/bin/honudb", "serve" ]