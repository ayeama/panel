#!/usr/bin/bash

if [ -f "/data/run.sh" ]; then
    exec /data/run.sh
fi

MANIFEST_URL="https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"
MANIFEST=$(curl -s $MANIFEST_URL)

VERSION="latest"
if [[ "$PANEL_VERSION" == "latest" ]]; then
    VERSION=$(echo $MANIFEST | jq -r '.latest.release')
else
    VERSION=$PANEL_VERSION
fi

if [ ! -f server.jar ]; then
    VERSION_LATEST_MANIFEST_URL=$(echo $MANIFEST | jq -r --arg version "${VERSION}" '.versions[] | select(.id==$version) | .url')
    VERSION_LATEST_MANIFEST=$(curl -s $VERSION_LATEST_MANIFEST_URL)
    VERSION_LATEST_SERVER_URL=$(echo $VERSION_LATEST_MANIFEST | jq -r '.downloads.server.url')
    wget -O server.jar $VERSION_LATEST_SERVER_URL
fi

if [[ ! -f eula.txt || ! -f server.properties ]]; then
    java -Xmx1024M -Xms1024M -jar server.jar --nogui --initSettings || exit 1
fi

if [ $PANEL_ACCEPT_EULA = true ]; then
    sed -i 's/^eula=false$/eula=true/' eula.txt
fi

exec java -Xmx1024M -Xms1024M -jar server.jar --nogui
