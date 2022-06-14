FROM golang:1.18

# Set timezone
ENV TZ Asia/Tokyo

# Update OS's packages
# Note that an error will occur if you separate
# the following commands by a backslash and a line break.
RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install -y vim

# Set the working directory
RUN mkdir -p /go/src/work
WORKDIR /go/src/work
ADD . /go/src/work

# Install packages for checking by static analysis
RUN go install github.com/kisielk/errcheck@latest \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
    go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest \
    curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
RUN cd /go/src/work/ops && go mod tidy
