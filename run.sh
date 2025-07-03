#!/usr/bin/bash

NAME=panel

podman pod create --name $NAME -p 8000:443
podman run -d --pod $NAME --name nginx --restart unless-stopped --volume ./nginx.conf:/etc/nginx/conf.d/panel.conf:z --volume ./.htpasswd:/.htpasswd:z --volume ./cert.pem:/cert.pem:z --volume ./cert.key:/cert.key:z docker.io/library/nginx:1.28.0-alpine
# podman run -d --pod $NAME --name ui --restart unless-stopped --security-opt label=disable localhost/panel/ui:latest
podman volume create ${NAME}_api_data
podman run -d --pod $NAME --name api --restart unless-stopped --env PANEL_SERVER_HOST=ayeama.com --env PANEL_RUNTIME_URI=unix:${XDG_RUNTIME_DIR}/podman/podman.sock --volume ${NAME}_api_data:/data --volume ${XDG_RUNTIME_DIR}/podman:${XDG_RUNTIME_DIR}/podman:Z --security-opt label=disable localhost/panel/api:latest
