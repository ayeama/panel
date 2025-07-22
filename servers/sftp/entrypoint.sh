#!/bin/ash

set -e

if [ -z "$PUBLIC_KEY" ]; then
    echo "PUBLIC_KEY is not set"
    exit 1
fi

if [ ! -f /etc/ssh/ssh_host_ed25519_key ]; then
  ssh-keygen -A
fi

mkdir -p /root/.ssh
echo "$PUBLIC_KEY" > /root/.ssh/authorized_keys

chmod 700 /root/.ssh
chmod 600 /root/.ssh/authorized_keys

exec /usr/sbin/sshd -D -e
