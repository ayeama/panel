#!/usr/bin/env bash

SAVE="saves/my-save.zip"

if [ ! -f "${SAVE}" ]; then
    bin/x64/factorio --create "${SAVE}"
fi

exec bin/x64/factorio --start-server "${SAVE}"
