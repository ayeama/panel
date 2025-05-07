#!/usr/bin/bash

set -x
set -e

NAME="My server"
WORLD="Dedicated"
PASSWORD="secret"

# install steamcmd
if [ ! -f /opt/steam/steamcmd.sh ]; then
    apt update -y
    apt-get install -y curl lib32gcc-s1

    mkdir /opt/steam
    cd /opt/steam

    curl -sqL "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz" | tar zxvf -
    ./steamcmd.sh +quit
fi

# install valheim
if [ ! -f /data/start_server.sh ]; then
    cd /data
    apt-get install -y libatomic1 libpulse-dev libpulse0
    /opt/steam/steamcmd.sh +force_install_dir /data +login anonymous +app_update 896660 validate +quit
fi

export templdpath=$LD_LIBRARY_PATH
export LD_LIBRARY_PATH=./linux64:$LD_LIBRARY_PATH
export SteamAppId=892970
./valheim_server.x86_64 -name "$NAME" -port 2456 -world "$WORLD" -password "$PASSWORD"
export LD_LIBRARY_PATH=$templdpath
