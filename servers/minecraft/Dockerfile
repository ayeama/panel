FROM docker.io/eclipse-temurin:21-jre-noble

WORKDIR /data
VOLUME ["/data"]

LABEL org.opencontainers.image.authors="github.com/ayeama/panel"
LABEL com.github.ayeama.panel.workload.name="Minecraft Vanilla"
LABEL com.github.ayeama.panel.workload.version="1"

ENV PANEL_VERSION="latest"
ENV PANEL_ACCEPT_EULA=true

RUN apt-get update -y && \
    apt-get install -y jq && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY ./entrypoint.sh /entrypoint.sh
CMD ["/entrypoint.sh"]

EXPOSE 25565
