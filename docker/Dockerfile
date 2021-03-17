# Simple usage with a mounted data directory:
# > docker build -t maticnetwork/heimdall:<tag> .
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.heimdalld:/root/.heimdalld maticnetwork/heimdall:<tag> heimdalld init

# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:latest as builder

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
RUN make build


# Pull all binaries into a second stage deploy alpine container
FROM ubuntu:latest

RUN apt update && apt install ca-certificates -y

# Copy required binarires to new container
COPY --from=builder /root/heimdall/build/* /usr/local/bin/


# add volumes
VOLUME ["/root/.heimdalld", "./logs" ]

# expose ports
EXPOSE 1317 26656 26657

# Run the binary.
CMD ["/usr/local/bin/heimdalld","start"]