package api

import "time"

type ServerResponse struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Image  string   `json:"image"`
	Status string   `json:"status"`
	Ports  []string `json:"ports"`
	// Env    map[string]string `json:"env"` // TODO implement
}

type ServerListResponse struct {
	Items []ServerResponse `json:"items"`
}

type ServerCreateRequest struct {
	Image string            `json:"image"`
	Env   map[string]string `json:"env"`
}

// NOTE backup manifest
// TODO add env to allow restore and import?
type ServerBackupManifest struct {
	Server  ServerBackupManifestServer `yaml:"server"`
	Created time.Time                  `yaml:"created"`
}

type ServerBackupManifestServer struct {
	Id     string                            `yaml:"id"`
	Name   string                            `yaml:"name"`
	Image  string                            `yaml:"image"`
	Mounts []ServerBackupManifestServerMount `yaml:"mounts"`
}

type ServerBackupManifestServerMount struct {
	Name        string `yaml:"name"`
	Destination string `yaml:"destination"`
}
