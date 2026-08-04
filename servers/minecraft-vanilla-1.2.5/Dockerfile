FROM docker.io/library/eclipse-temurin:25-jre-alpine

LABEL com.github.ayeama.panel.image.id="d4114e4a-8de0-4719-8c6a-d45346f57009"
LABEL com.github.ayeama.panel.image.version="1"
LABEL com.github.ayeama.panel.image.name="Minecraft"
LABEL com.github.ayeama.panel.image.variant="Vanilla"
LABEL com.github.ayeama.panel.image.release="1.2.5"
LABEL com.github.ayeama.panel.image.description=""

WORKDIR /app

RUN wget https://launcher.mojang.com/v1/objects/d8321edc9470e56b8ad5c67bbd16beba25843336/server.jar

COPY docker-entrypoint.sh /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]

EXPOSE 25565/tcp

VOLUME /app
