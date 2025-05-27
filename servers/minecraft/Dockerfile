# TODO: https://github.com/krallin/tini

FROM docker.io/eclipse-temurin:21-jre-noble

STOPSIGNAL SIGTERM

WORKDIR /data
VOLUME ["/data"]

RUN apt-get update -y && apt-get install -y jq

COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]

EXPOSE 25565
