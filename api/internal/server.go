package internal

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/handler"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/ayeama/panel/api/internal/runtime"
	"github.com/ayeama/panel/api/internal/service"
)

type Server struct {
	server http.Server
}

func NewServer() *Server {
	db, err := sql.Open("sqlite3", "panel.db")
	if err != nil {
		panic(err)
	}

	runtime, err := runtime.New(runtime.RuntimeTypePodman)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	nodeRepository := repository.NewNodeRepository(db)
	nodeService := service.NewNodeService(nodeRepository)
	nodeHandler := handler.NewNodeHandler(nodeService)
	nodeHandler.RegisterHandlers(mux)

	manifestRepository := repository.NewManifestRepository(db)
	manifestService := service.NewManifestService(manifestRepository)
	manifestHandler := handler.NewManifestHandler(manifestService)
	manifestHandler.RegisterHandlers(mux)

	serverRepository := repository.NewServerRepository(db)
	serverService := service.NewServerService(runtime, serverRepository, nodeRepository, manifestRepository)
	serverHandler := handler.NewServerHandler(serverService)
	serverHandler.RegisterHandlers(mux)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	userHandler.RegisterHandlers(mux)

	nodesPaginated := nodeService.Read(domain.Pagination{Limit: 20, Offset: 0})
	for _, node := range nodesPaginated.Items {
		fmt.Println(node)
		runtime.AddNode(&node)
	}

	server := Server{
		server: http.Server{
			Addr:    "0.0.0.0:8000",
			Handler: Log(Cors(mux)), // TODO update middleware method
		},
	}

	return &server
}

func (s *Server) Start() {
	slog.Info("starting")
	// TODO run in go routine?
	s.server.ListenAndServeTLS("cert.pem", "key.pem")
}
