FROM docker.io/library/eclipse-temurin:25-jre-alpine

LABEL com.github.ayeama.panel.image.id="d4114e4a-8de0-4719-8c6a-d45346f57009"
LABEL com.github.ayeama.panel.image.version="1"
LABEL com.github.ayeama.panel.image.name="Minecraft"
LABEL com.github.ayeama.panel.image.variant="Vanilla"
LABEL com.github.ayeama.panel.image.release="26.2"
LABEL com.github.ayeama.panel.image.description=""

WORKDIR /app

COPY docker-entrypoint.sh /docker-entrypoint.sh

RUN wget https://piston-data.mojang.com/v1/objects/823e2250d24b3ddac457a60c92a6a941943fcd6a/server.jar

ENTRYPOINT ["/docker-entrypoint.sh"]

EXPOSE 25565/tcp

VOLUME /app
