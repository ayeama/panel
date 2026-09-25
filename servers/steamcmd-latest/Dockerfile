FROM docker.io/library/debian:trixie-slim

RUN apt-get update && \
    apt-get install --no-install-recommends -y \
        curl \
        ca-certificates \
        lib32gcc-s1 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /opt/steam && \
    cd /opt/steam && \
    curl -sqL "https://client-update.steamstatic.com/installer/steamcmd_linux.tar.gz" | tar zxvf - && \
    ./steamcmd.sh +quit
