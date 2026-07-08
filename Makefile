.PHONY: cert run build deploy clean

run:
	go run cmd/panel/main.go

build:
	podman build -f Dockerfile.backend -t panel/backend:0.0.1 .
	podman build -f Dockerfile.frontend -t panel/frontend:0.0.1 .

deploy:
	podman pod create --name panel -p 8080:8080 --userns keep-id
	podman run --pod panel --name backend -d --security-opt label=disable -v "/run/user/1000/podman/podman.sock:/run/user/1000/podman/podman.sock:Z" panel/backend:0.0.1
	podman run --pod panel --name frontend -d panel/frontend:0.0.1

clean:
	podman pod rm -f panel
