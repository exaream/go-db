FROM golang:1.19-alpine3.16

# Set timezone.
ENV TZ Asia/Tokyo

# Update Alpine Linux.
RUN apk update && \
    apk upgrade && \
    apk add alpine-sdk build-base vim

# Add a group and a user for local environment.
# Do NOT use the following setting for production environment.
RUN addgroup -g 1000 samplegroup && \
    adduser -S -D -u 1000 -G samplegroup sampleuser

# Set the working directory.
RUN mkdir -p /go/src/work
WORKDIR /go/src/work
ADD . /go/src/work
RUN chown -R sampleuser:samplegroup /go/
USER sampleuser

# Install Go packages.
RUN go install golang.org/x/vuln/cmd/govulncheck@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin && \
    cd /go/src/work && \
    go mod tidy
