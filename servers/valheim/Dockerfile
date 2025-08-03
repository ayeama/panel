FROM docker.io/library/debian:12.11

LABEL org.opencontainers.image.authors="github.com/ayeama/panel"
LABEL com.github.ayeama.panel.workload.name="Valheim"
LABEL com.github.ayeama.panel.workload.version="1"
LABEL com.github.ayeama.panel.workload.stop="stop"

STOPSIGNAL SIGTERM

WORKDIR /data
VOLUME ["/data"]

RUN apt-get update -y && \
    apt-get install -y \
        curl \
        lib32gcc-s1 \
        libatomic1 \
        libpulse-dev \
        libpulse0

COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]

EXPOSE 2456
