FROM golang:1.18

# Set timezone
ENV TZ Asia/Tokyo

# Update OS's packages
RUN apt-get update

# Set the working directory
RUN mkdir /go/src/work
WORKDIR /go/src/work
ADD . /go/src/work

# Install packages for checking by static analysis
RUN go install github.com/kisielk/errcheck@latest \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
    go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest \
    curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
