FROM golang:latest

ARG HEIMDALL_DIR=/var/lib/heimdall
ENV HEIMDALL_DIR=$HEIMDALL_DIR

RUN apt-get update -y && apt-get upgrade -y \
    && apt install build-essential git -y \
    && mkdir -p $HEIMDALL_DIR

WORKDIR ${HEIMDALL_DIR}
COPY . .

RUN make install
RUN groupadd -g 20137 heimdall \
    && useradd -u 20137 --no-log-init --create-home -r -g heimdall heimdall \
    && chown -R heimdall:heimdall ${HEIMDALL_DIR}

COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh

USER heimdall
ENV SHELL /bin/bash
EXPOSE 1317 26656 26657

ENTRYPOINT ["entrypoint.sh"]
