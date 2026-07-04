CERT := server.crt
KEY := server.key

.PHONY: cert run clean

cert: $(CERT) $(KEY)

$(CERT) $(KEY):
	openssl req -x509 -newkey rsa:4096 -sha256 -days 365 -nodes -keyout $(KEY) -out $(CERT) -subj "/CN=localhost" -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

run:
	go run cmd/panel/main.go

clean:
	rm -rf $(CERT) $(KEY)
