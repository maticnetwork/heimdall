FROM golang:latest AS builder

ARG HEIMDALL_DIR=/heimdall
ENV HEIMDALL_DIR=$HEIMDALL_DIR

RUN apt-get update -y && apt-get upgrade -y \
    && apt install build-essential git -y \
    && mkdir -p /heimdall

WORKDIR ${HEIMDALL_DIR}
COPY . .

RUN make install

# Seconds stage
FROM alpine:3.15

WORKDIR /app
COPY --from=builder /go/bin/* /app/

RUN apk add --no-cache --virtual=.build-dependencies wget ca-certificates \
    && wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
    && wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.34-r0/glibc-2.34-r0.apk \
    && apk add --no-cache glibc-2.34-r0.apk \
    && rm glibc-2.34-r0.apk \
    && apk del .build-dependencies

ENV PATH /app:$PATH

# add volumes
VOLUME [ "/root/.heimdalld" ]

# expose ports
EXPOSE 1317 26656 26657
