#!/usr/bin/env bash

if [ ! -f config.txt ]; then
    cat > config.txt <<EOF
autocreate=3
banlist="banlist.txt"
difficulty=0
language="en-US"
max_players="8"
motd=""
npcstream="60"
password="secret"
port="7777"
priority="1"
secure="1"
seed=
upnp="0"
world_name="world"
world_path="/data/worlds/"
world="/data/worlds/world.wld"
EOF
fi

exec ./TerrariaServer.bin.x86_64 -config config.txt
