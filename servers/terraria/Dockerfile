FROM docker.io/library/debian:12.11

LABEL org.opencontainers.image.authors="github.com/ayeama/panel"
LABEL com.github.ayeama.panel.workload.name="Terraria"
LABEL com.github.ayeama.panel.workload.version="1"
LABEL com.github.ayeama.panel.workload.stop="exit"

ENV PANEL_VERSION="1449"
ENV PANEL_WORLD_SIZE="1"
ENV PANEL_SEED=""
ENV PANEL_DIFFICULTY="0"
ENV PANEL_PASSWORD=""

WORKDIR /data
VOLUME ["/data"]

RUN apt-get update -y && \
    apt-get install -y \
        wget \
        unzip && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]

EXPOSE 7777
