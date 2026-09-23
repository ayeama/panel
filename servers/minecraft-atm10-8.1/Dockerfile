FROM docker.io/library/eclipse-temurin:25-jre-alpine

LABEL com.github.ayeama.panel.image.id="b2ea2b3c-1d33-4146-901a-111f0e975890"
LABEL com.github.ayeama.panel.image.version="1"
LABEL com.github.ayeama.panel.image.name="Minecraft"
LABEL com.github.ayeama.panel.image.variant="All The Mods 10"
LABEL com.github.ayeama.panel.image.release="8.1"
LABEL com.github.ayeama.panel.image.description=""

WORKDIR /app

RUN wget -O ServerFiles-8.1.zip https://www.curseforge.com/api/v1/mods/925200/files/8764245/download && \
    unzip ServerFiles-8.1.zip && \
    rm -rf ServerFiles-8.1.zip && \
    wget -O neoforge-21.1.249-installer.jar https://maven.neoforged.net/releases/net/neoforged/neoforge/21.1.249/neoforge-21.1.249-installer.jar && \
    java -jar neoforge-21.1.249-installer.jar -installServer

COPY docker-entrypoint.sh /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]

EXPOSE 25565/tcp

VOLUME /app
