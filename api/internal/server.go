package internal

import (
	"database/sql"
	"log/slog"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/errgroup"

	"github.com/ayeama/panel/api/internal/config"
	"github.com/ayeama/panel/api/internal/handler"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/ayeama/panel/api/internal/runtime"
	"github.com/ayeama/panel/api/internal/service"
)

type Server struct {
	server        http.Server
	serverService *service.ServerService
}

func NewServer() *Server {
	db, err := sql.Open("sqlite3", "/data/panel.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS images (image VARCHAR(255) PRIMARY KEY)")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS servers (id VARCHAR(36) PRIMARY KEY,name VARCHAR(256) NOT NULL,status VARCHAR(32) NOT NULL,container_id VARCHAR(64) NOT NULL,container_port VARCHAR(5))")
	if err != nil {
		panic(err)
	}

	runtime, err := runtime.New()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	imageRepository := repository.NewImageRepository(db)
	imageService := service.NewImageService(imageRepository)
	imageHandler := handler.NewImageHandler(imageService)
	imageHandler.RegisterHandlers(mux)

	serverRepository := repository.NewServerRepository(db)
	serverService := service.NewServerService(runtime, serverRepository, imageService)
	serverHandler := handler.NewServerHandler(serverService)
	serverHandler.RegisterHandlers(mux)

	eventHandler := handler.NewEventHandler(serverService)
	eventHandler.RegisterHandlers(mux)

	server := Server{
		server: http.Server{
			Addr:    config.Config.ApiAddress,
			Handler: Log(Cors(mux)), // TODO update middleware method
		},
		serverService: serverService,
	}

	return &server
}

func (s *Server) Start() {
	slog.Info("starting", slog.String("address", s.server.Addr))

	var eg errgroup.Group

	api := func() error {
		return s.server.ListenAndServeTLS("cert.pem", "key.pem")
	}

	events := func() error {
		return s.serverService.Events()
	}

	eg.Go(api)
	eg.Go(events)

	err := eg.Wait()
	if err != nil {
		panic(err)
	}
}
