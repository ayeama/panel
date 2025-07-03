# Panel

## Requirements

* Disk size quotas require the XFS filesystem.


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


# UML

## servers

| id (PK) | name | status |
| --- | --- | --- |
| 29def84d-46b4-43f6-b3c7-a5852f132908 | Minecraft Vanilla | running |

## containers

| id (PK) | name | status | server_id (FK) |
| --- | --- | --- | --- |
| a44c2123761f9a30aecb863a3f35caf0a125d7b3264ecbe52aba031dbe5db110 | lucid_boyd | running | 29def84d-46b4-43f6-b3c7-a5852f132908 |

## users

| id (PK) | username | password |
| --- | --- | --- |
| 4f56915f-f087-4100-8d9c-ff1b2c2b5440 | ayeama | admin |
| b29331c7-57a3-4b81-815a-5e2efa193529 | john | smith |

## roles

| id (PK) | name |
| --- | --- |
| 1 | owner |
| 2 | member |

## user_server_roles

| user_id (PK, FK) | server_id (PK, FK) | role_id (FK) |
| --- | --- | --- |
| 4f56915f-f087-4100-8d9c-ff1b2c2b5440 | 29def84d-46b4-43f6-b3c7-a5852f132908 | 1 |

## invite

| id (PK) | server_id (FK) | inviter_user_id (FK) | invited_user_id (FK) | role_id | status | invited_date | accepted_date |
| --- | --- | --- | --- | --- | --- | --- | --- |


# DNS

```sh
sudo dnf install -y podman-plugins
```
