package internal

import (
	"database/sql"
	"log/slog"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/ayeama/panel/api/internal/runtime"
	"github.com/ayeama/panel/api/internal/service"
)

type Agent struct {
	serverService *service.ServerService
}

func NewAgent() *Agent {
	db, err := sql.Open("sqlite3", "panel.db")
	if err != nil {
		panic(err)
	}

	runtime, err := runtime.New(runtime.RuntimeTypePodman)
	if err != nil {
		panic(err)
	}

	nodeRepository := repository.NewNodeRepository(db)
	manifestRepository := repository.NewManifestRepository(db)

	serverRepository := repository.NewServerRepository(db)
	serverService := service.NewServerService(runtime, serverRepository, nodeRepository, manifestRepository)

	return &Agent{serverService: serverService}
}

func (a *Agent) Start() {
	slog.Info("starting")

	for event := range a.serverService.Events() {
		switch e := event.(type) {
		case domain.ServerCreatedEvent:
			server := domain.Server{Id: e.Id, Status: "created"}
			a.serverService.Update(&server)
		case domain.ServerStartedEvent:
			server := domain.Server{Id: e.Id, Status: "running"}
			a.serverService.Update(&server)
		case domain.ServerStoppedEvent:
			server := domain.Server{Id: e.Id, Status: "stopped"}
			a.serverService.Update(&server)
		// case domain.ServerDeletedEvent:
		// 	server := domain.Server{Id: e.Id, Status: "created"}
		// 	server.Container = &domain.Container{Id: e.ContainerId}
		// 	a.serverService.Update(&server)
		default:
		}
	}
}
