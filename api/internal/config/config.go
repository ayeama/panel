package config

import (
	"os"
)

var Config *config

type config struct {
	ApiAddress      string
	ServerHost      string
	ServerPortRange string
	RuntimeUri      string
}

func New() {
	apiAddress := os.Getenv("PANEL_ADDRESS")
	if apiAddress == "" {
		apiAddress = "0.0.0.0:8000"
	}

	serverHost := os.Getenv("PANEL_SERVER_HOST")
	if serverHost == "" {
		serverHost = "127.0.0.1"
		// panic("config: server host requried")
	}

	serverPortRange := os.Getenv("PANEL_SERVER_PORT_RANGE")
	if serverPortRange == "" {
		serverPortRange = "45000-45099"
	}

	runtimeUri := os.Getenv("PANEL_RUNTIME_URI")
	if runtimeUri == "" {
		runtimeUri = "unix:/run/user/1000/podman/podman.sock"
	}

	Config = &config{
		ApiAddress:      apiAddress,
		ServerHost:      serverHost,
		ServerPortRange: serverPortRange,
		RuntimeUri:      runtimeUri,
	}
}
