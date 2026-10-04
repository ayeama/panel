# https://www.valheimgame.com/support/a-guide-to-dedicated-servers/

FROM panel/steamcmd:1

LABEL com.github.ayeama.panel.image.id="a80f899f-11e6-4b9c-b282-6740ede44adf"
LABEL com.github.ayeama.panel.image.version="1"
LABEL com.github.ayeama.panel.image.name="Valheim"
LABEL com.github.ayeama.panel.image.variant=""
LABEL com.github.ayeama.panel.image.release="latest"
LABEL com.github.ayeama.panel.image.description=""

WORKDIR /app

RUN apt-get update && \
    apt-get install --no-install-recommends -y \
        libpulse-dev \
        libatomic1 \
        libpulse0 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    /opt/steam/steamcmd.sh +@sSteamCmdForcePlatformType linux +force_install_dir /app +login anonymous +app_update 896660 -beta public validate +quit && \
    cp docker/steamclient.so /usr/lib

COPY docker-entrypoint.sh /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]

EXPOSE 2456/udp

VOLUME /app
