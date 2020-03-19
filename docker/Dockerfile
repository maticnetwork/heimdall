# Simple usage with a mounted data directory:
# > docker build -t heimdall .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.heimdalld:/root/.heimdalld heimdall heimdalld init

# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:latest

# update available packages
RUN apt-get update -y && apt-get upgrade -y && apt install build-essential -y

# setup dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# create go src directory and clone heimdall
RUN mkdir -p /go/src/github.com/maticnetwork/heimdall \
  && cd /go/src/github.com/maticnetwork/heimdall

ADD . /go/src/github.com/maticnetwork/heimdall/

# change work directory
WORKDIR /go/src/github.com/maticnetwork/heimdall

# GOBIN required for go install
ENV GOBIN $GOPATH/bin

# run build
RUN make install

# add volumes
VOLUME [ "/root/.heimdalld", "./logs" ]

# expose ports
EXPOSE 1317 26656 26657