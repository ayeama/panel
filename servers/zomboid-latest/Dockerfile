# https://pzwiki.net/wiki/Dedicated_server#Linux
# 4 cpu 16 memory

FROM panel/steamcmd:1

LABEL com.github.ayeama.panel.image.id="703bb78a-29ac-401c-8360-a5d160fd2d61"
LABEL com.github.ayeama.panel.image.version="1"
LABEL com.github.ayeama.panel.image.name="Zomboid"
LABEL com.github.ayeama.panel.image.variant=""
LABEL com.github.ayeama.panel.image.release="latest"
LABEL com.github.ayeama.panel.image.description=""

WORKDIR /app

RUN /opt/steam/steamcmd.sh +@sSteamCmdForcePlatformType linux +force_install_dir /app +login anonymous +app_update 380870 validate +quit

COPY docker-entrypoint.sh /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]

EXPOSE 16261/udp

VOLUME /app
