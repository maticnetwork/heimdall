# Simple usage with a mounted data directory:
# > docker build -t maticnetwork/heimdall:<tag> .
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.heimdalld:/root/.heimdalld maticnetwork/heimdall:<tag> heimdalld init

# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:latest

# update available packages
RUN apt-get update -y && apt-get upgrade -y && apt install build-essential -y

# create go src directory and clone heimdall
RUN mkdir -p /root/heimdall

# add code to docker instance
ADD . /root/heimdall/

# change work directory
WORKDIR /root/heimdall

# GOBIN required for go install
ENV GOBIN $GOPATH/bin

# run build
RUN make install

# add volumes
VOLUME [ "/root/.heimdalld", "./logs" ]

# expose ports
EXPOSE 1317 26656 26657
