FROM docker.io/library/debian:12.11

WORKDIR /data
VOLUME ["/data"]

RUN apt-get update -y && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

EXPOSE 8000

LABEL org.opencontainers.image.authors="github.com/ayeama/panel"
LABEL com.github.ayeama.panel.workload.name="noop"
LABEL com.github.ayeama.panel.workload.version="1"
LABEL com.github.ayeama.panel.workload.stop=""

ENV PANEL_SECONDS="86400"
