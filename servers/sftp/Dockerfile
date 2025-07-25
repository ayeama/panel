FROM docker.io/alpine:3.22.0

WORKDIR /data
VOLUME ["/data"]

LABEL org.opencontainers.image.authors="github.com/ayeama/panel"
LABEL com.github.ayeama.panel.name="SFTP"
LABEL com.github.ayeama.panel.version="1"

ENV PUBLIC_KEY=

RUN apk add --no-cache openssh

COPY ./entrypoint.sh /entrypoint.sh
COPY ./add-key /usr/local/bin/add-key
COPY ./sshd_config /etc/ssh/sshd_config

CMD ["/entrypoint.sh"]

EXPOSE 22/tcp
