# https://wiki.vintagestory.at/Guide:Dedicated_Server#Dedicated_server_on_Linux
# 2 cpu 2 memory

FROM mcr.microsoft.com/dotnet/runtime:10.0-noble

LABEL com.github.ayeama.panel.image.id="12c3e237-18e2-47f5-bf20-0951ef585321"
LABEL com.github.ayeama.panel.image.version="1"
LABEL com.github.ayeama.panel.image.name="Vintage Story"
LABEL com.github.ayeama.panel.image.variant=""
LABEL com.github.ayeama.panel.image.release="1.22.6"
LABEL com.github.ayeama.panel.image.description=""

WORKDIR /app

RUN apt-get update && \
    apt-get install -y \
        wget && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    wget https://cdn.vintagestory.at/gamefiles/stable/vs_server_linux-x64_1.22.6.tar.gz && \
    tar -zxf vs_server_linux-x64_1.22.6.tar.gz && \
    rm -rf vs_server_linux-x64_1.22.6.tar.gz

COPY docker-entrypoint.sh /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]

EXPOSE 42420

VOLUME /app
