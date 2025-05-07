#!/usr/bin/bash

set -x
set -e

MANIFEST_URL="https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"
ACCEPT_EULA=true

if [ ! -f server.jar ]; then
    MANIFEST=$(curl -s $MANIFEST_URL)
    VERSION_LATEST=$(echo $MANIFEST | jq -r '.latest.release')
    VERSION_LATEST_MANIFEST_URL=$(echo $MANIFEST | jq -r --arg version "${VERSION_LATEST}" '.versions[] | select(.id==$version) | .url')
    VERSION_LATEST_MANIFEST=$(curl -s $VERSION_LATEST_MANIFEST_URL)
    VERSION_LATEST_SERVER_URL=$(echo $VERSION_LATEST_MANIFEST | jq -r '.downloads.server.url')
    wget -O server.jar $VERSION_LATEST_SERVER_URL
fi

if [[ ! -f eula.txt || ! -f server.properties ]]; then
    java -Xmx1024M -Xms1024M -jar server.jar --nogui --initSettings || exit 1
fi

if [ $ACCEPT_EULA = true ]; then
    sed -i 's/^eula=false$/eula=true/' eula.txt
fi

exec java -Xmx1024M -Xms1024M -jar server.jar --nogui

