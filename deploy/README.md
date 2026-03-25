# Installation

```sh
sudo useradd panel
sudo passwd panel
sudo usermod -aG wheel panel

sudo dnf update -y
sudo dnf install -y podman

systemctl --user start podman.socket
loginctl enable-linger $USER

sudo firewall-cmd --add-port=8442/tcp --add-port=8443/tcp --permanent
sudo firewall-cmd --reload
sudo firewall-cmd --list-all
```

```sh
ansible-playbook -i inventory.yml -J playbook.yml
```

Reserve ephemeral ports for exposing servers.

```sh
echo "net.ipv4.ip_local_reserved_ports = 45000-46023" | sudo tee -a /etc/sysctl.d/99-panel-ports.conf
sudo sysctl --system
```

```sh
sudo firewall-cmd --add-port=45000-46023/tcp --add-port=45000-46023/udp --permanent
sudo firewall-cmd --reload
```
