#!/usr/bin/bash

if [ -f /data/run.sh ]; then
    exec /data/run.sh
fi

# PANEL_VERSION

WORLD="/data/worlds/world.wld"
WORLD_SIZE="1"
SEED=""
WORLD_NAME="world"
DIFFICULTY="0"
MAX_PLAYERS="8"
PORT="7777"
PASSWORD=""
MOTD=""
WORLD_PATH="/data/worlds/"
BANLIST="banlist.txt"
SECURE="1"
LANGUAGE="en-US"
UPNP="0"
NPCSTREAM="60"
PRIORITY="1"

if [ ! -f terraria-server.zip ]; then
    wget -O terraria-server.zip "https://terraria.org/api/download/pc-dedicated-server/terraria-server-${PANEL_VERSION}.zip"
fi

if [ ! -f TerrariaServer.bin.x86_64 ]; then
    unzip terraria-server.zip
    mv ${PANEL_VERSION}/Linux/* .
    rm -rf ${PANEL_VERSION}/ # TODO delete? terraria-server.zip
    chmod +x TerrariaServer.bin.x86_64
fi

mkdir -p worlds

if [ ! -f config.txt ]; then
    cat > config.txt <<EOF
world=${WORLD}
autocreate=${WORLD_SIZE}
seed=${SEED}
worldname=${WORLD_NAME}
difficulty=${DIFFICULTY}
maxplayers=${MAX_PLAYERS}
port=${PORT}
password=${PASSWORD}
motd=${MOTD}
worldpath=${WORLD_PATH}
banlist=${BANLIST}
secure=${SECURE}
language=${LANGUAGE}
upnp=${UPNP}
npcstream=${NPCSTREAM}
priority=${PRIORITY}
EOF
fi

exec ./TerrariaServer.bin.x86_64 -config config.txt
