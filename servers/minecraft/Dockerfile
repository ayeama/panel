FROM docker.io/eclipse-temurin:25-jre-noble

LABEL org.opencontainers.image.authors="github.com/ayeama/panel"
LABEL com.github.ayeama.panel.server.name="Minecraft Vanilla"
LABEL com.github.ayeama.panel.server.version="1"
LABEL com.github.ayeama.panel.server.stop="stop"

ENV PANEL_VERSION="latest"
ENV PANEL_ACCEPT_EULA=true

WORKDIR /data
VOLUME ["/data"]

RUN apt-get update -y && \
    apt-get install -y \
        curl \
        jq && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

EXPOSE 25565
