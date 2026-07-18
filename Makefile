.PHONY: cert run build deploy clean

run:
	go run cmd/panel/main.go

build:
	podman build -f Dockerfile.backend --ignorefile .dockerignore.backend --target backend -t panel/backend:0.0.1 .
	podman build -f Dockerfile.backend --ignorefile .dockerignore.backend --target backend-dns -t panel/backend-dns:0.0.1 .
	podman build -f Dockerfile.backend --ignorefile .dockerignore.backend --target backend-notify -t panel/backend-notify:0.0.1 .
	podman build -f Dockerfile.frontend --ignorefile .dockerignore.frontend -t panel/frontend:0.0.1 .

deploy:
	podman pod create --name panel -p 8080:8080 --userns keep-id
	podman run --pod panel --name backend --restart unless-stopped -d --security-opt label=disable -v "/run/user/1000/podman/podman.sock:/run/user/1000/podman/podman.sock:Z" panel/backend:0.0.1
	podman run --pod panel --name backend-dns --restart unless-stopped -e "PANEL_DNS_HOST" -e "CLOUDFLARE_ZONE_ID" -e "CLOUDFLARE_API_TOKEN" -d panel/backend-dns:0.0.1
	podman run --pod panel --name backend-notify --restart unless-stopped -e "PANEL_NOTIFY_NTFY_HOST" -e "PANEL_NOTIFY_NTFY_TOPIC" -e "PANEL_NOTIFY_NTFY_TOKEN" -d panel/backend-notify:0.0.1
	podman run --pod panel --name frontend --restart unless-stopped -d panel/frontend:0.0.1

clean:
	podman pod rm -f panel
