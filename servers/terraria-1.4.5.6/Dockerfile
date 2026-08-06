# https://terraria.wiki.gg/wiki/Server

FROM docker.io/library/debian:trixie-slim

LABEL com.github.ayeama.panel.image.id="7077c690-a4c7-4f67-83fd-b1b697d370d4"
LABEL com.github.ayeama.panel.image.version="1"
LABEL com.github.ayeama.panel.image.name="Terraria"
LABEL com.github.ayeama.panel.image.variant=""
LABEL com.github.ayeama.panel.image.release="1.4.5.6"
LABEL com.github.ayeama.panel.image.description=""

WORKDIR /app

RUN apt-get update && \
    apt-get install --no-install-recommends -y \
        wget \
        ca-certificates \
        unzip && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    wget https://terraria.org/api/download/pc-dedicated-server/terraria-server-1456.zip && \
    unzip terraria-server-1456.zip && \
    mv /app/1456/Linux/* /app/ && \
    chmod +x TerrariaServer.bin.x86_64 && \
    rm -rf 1456 terraria-server-1456.zip

COPY docker-entrypoint.sh /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]

EXPOSE 7777/tcp

VOLUME /app
