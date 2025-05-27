FROM docker.io/library/debian:12.11

STOPSIGNAL SIGTERM

WORKDIR /data
VOLUME ["/data"]

RUN apt-get update -y

COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]

EXPOSE 2456
