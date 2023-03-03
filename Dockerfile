FROM golang:1.19.1-bullseye AS builder

RUN apt-get update && \
    apt-get install -y --no-install-recommends curl libsodium-dev && \
    curl https://download.joinself.com/olm/libself-olm_0.1.17_amd64.deb -o /tmp/libself-olm_0.1.17_amd64.deb && \
    curl https://download.joinself.com/omemo/libself-omemo_0.1.3_amd64.deb -o /tmp/libself-omemo_0.1.3_amd64.deb && \
    apt-get install -y --no-install-recommends /tmp/libself-olm_0.1.17_amd64.deb /tmp/libself-omemo_0.1.3_amd64.deb && \
    rm -rf /var/lib/apt/lists/* /tmp/*
