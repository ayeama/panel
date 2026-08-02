# https://wiki.factorio.com/Multiplayer#Dedicated/Headless_server

FROM docker.io/library/debian:trixie-slim

LABEL com.github.ayeama.panel.image.id="9fc73bfb-31d2-40cd-80f4-c98e2ab931f0"
LABEL com.github.ayeama.panel.image.version="1"
LABEL com.github.ayeama.panel.image.name="Factorio"
LABEL com.github.ayeama.panel.image.variant=""
LABEL com.github.ayeama.panel.image.release="2.0.77"
LABEL com.github.ayeama.panel.image.description=""

WORKDIR /app

COPY docker-entrypoint.sh /docker-entrypoint.sh

RUN apt-get update && \
    apt-get install -y \
        curl \
        xz-utils && \
    apt-get clean

RUN curl -L -o linux64 https://www.factorio.com/get-download/2.0.77/headless/linux64 && \
    tar --strip-components=1 -xf linux64 && \
    rm -rf linux64

ENTRYPOINT ["/docker-entrypoint.sh"]

EXPOSE 34197/udp

VOLUME /app
