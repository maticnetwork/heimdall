# Simple usage with a mounted data directory:
# > docker build -t maticnetwork/heimdall:<tag> .
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.heimdalld:/root/.heimdalld maticnetwork/heimdall:<tag> heimdalld init

# Build Geth in a stock Go builder container
FROM golang:latest as builder

# RUN apk add --no-cache make gcc musl-dev linux-headers git bash
RUN apt-get update -y && apt-get upgrade -y && apt install build-essential -y

ADD . /root/heimdall
RUN cd /root/heimdall && make build

CMD ["/bin/bash"]

# Pull heimdall into a second stage deploy alpine container
FROM alpine:latest

COPY --from=builder /root/heimdall/build/* /usr/local/bin/

# add volumes
VOLUME [ "/root/heimdall" ]

EXPOSE 1317 26656 26657
