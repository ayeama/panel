#!/usr/bin/env bash

NAME="My server"
PASSWORD="secret"
WORLD="Dedicated"
SAVEDIR="/app/Valheim"

export templdpath=$LD_LIBRARY_PATH
export LD_LIBRARY_PATH=/app/linux64:$LD_LIBRARY_PATH
export SteamAppId=892970
exec /app/valheim_server.x86_64 -name "$NAME" -port 2456 -world "$WORLD" -password "$PASSWORD" -savedir "$SAVEDIR"
