# Panel

systemctl --user start podman.socket


podman run --rm -it --cpus 1 --memory 1G -v "./data:/data" -p 25565:25565 eclipse-temurin:21-jre bash

cd /data
java -Xmx1024M -Xms1024M -jar server.jar nogui


podman run --rm --cpus 1 --memory 1G -v "./data:/data" -p 25565:25565 panel:latest

podman run -di --cpus 1 --memory 1G -v "./data:/data" -p 25565:25565 panel:latest
podman run -di --cpus 1 --memory 1G -p 25565:25565 panel:latest

```sh
openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out cert.pem -days 365 -subj "/CN=localhost"
```

If using self signed certs for development, before running the UI manually vist a API endpoint and trust the certificate, otherwise it will not work.


# Start service

```sh
systemctl --user start podman.socket

podman system service -t 0
```

# Dependencies

```sh
sudo dnf -y install podman
sudo dnf -y install catatonit conmon containers-common-extra

sudo dnf instal -y btrfs-progs-devel gpgme-devel
```


# Frontend

vue with bootstrap

## Dependencies

```sh
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | bash
nvm install --lts
nvm use --lts
```

# Steam Games

[steamdb](https://steamdb.info) for possible manifests etc?
