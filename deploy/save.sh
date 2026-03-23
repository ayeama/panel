#!/usr/bin/env bash

# docker compose build --parallel

docker save -o api.tar panel-backend:latest
docker save -o ui.tar panel-frontend:latest

scp api.tar panel:
scp ui.tar panel:

rm -rf api.tar ui.tar

# podman save -o minecraft.tar minecraft:0.0.1-jre21
# podman save -o valheim.tar valheim:0.0.1
# podman save -o terraria.tar terraria:0.0.1
