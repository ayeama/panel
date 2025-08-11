#!/usr/bin/bash

if [ -f "/data/run.sh" ]; then
    exec /data/run.sh
fi

SECONDS=$PANEL_SECONDS

echo "sleeping for ${SECONDS} seconds..."
sleep "${SECONDS}s"
