package domain

type Container struct {
	Id     string
	Name   string
	Status string
}

const (
	ServerStatusCreating string = "creating"
	ServerStatusCreated  string = "created"
	ServerStatusStarted  string = "started"
	ServerStatusRunning  string = "running"
	ServerStatusStopping string = "stopping"
	ServerStatusStopped  string = "stopped"
	ServerStatusDeleting string = "deleting"
)

type Server struct {
	Id        string
	Name      string
	Status    string
	Container *Container
}
