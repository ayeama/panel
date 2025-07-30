FROM docker.io/library/debian:12.11

WORKDIR /data
VOLUME ["/data"]

LABEL org.opencontainers.image.authors="github.com/ayeama/panel"
LABEL com.github.ayeama.panel.workload.name="Terraria"
LABEL com.github.ayeama.panel.workload.version="1"

ENV PANEL_VERSION="1449"

RUN apt-get update -y && \
    apt-get install -y \
        wget \
        unzip

COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]

EXPOSE 7777
