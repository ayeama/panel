#!/usr/bin/bash
# https://www.valheimgame.com/support/a-guide-to-dedicated-servers/

if [ -f /data/run.sh ]; then
    exec /data/run.sh
fi

NAME=$PANEL_NAME
PASSWORD=$PANEL_PASSWORD

WORLD="Dedicated"
SAVEDIR="/data/Valheim"

mkdir -p $SAVEDIR

# install steamcmd
if [ ! -f /opt/steam/steamcmd.sh ]; then
    mkdir /opt/steam
    cd /opt/steam

    curl -sqL "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz" | tar zxvf -
    ./steamcmd.sh +quit
fi

if [ ! -f /data/start_server.sh ]; then
    cd /data
    /opt/steam/steamcmd.sh +force_install_dir /data +login anonymous +app_update 896660 validate +quit
fi

export templdpath=$LD_LIBRARY_PATH
export LD_LIBRARY_PATH=./linux64:$LD_LIBRARY_PATH
export SteamAppId=892970

exec ./valheim_server.x86_64 -name "$NAME" -port 2456 -world "$WORLD" -password "$PASSWORD" -savedir "$SAVEDIR"
