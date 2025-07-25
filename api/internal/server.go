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
	db, err := sql.Open("sqlite3", "file:panel.db?_foreign_keys=true&_busy_timeout=1000000") // TODO
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS images(
			id TEXT NOT NULL UNIQUE PRIMARY KEY,
			tag TEXT NOT NULL UNIQUE
		);
		CREATE TABLE IF NOT EXISTS servers(
			id TEXT NOT NULL UNIQUE PRIMARY KEY,
			image_id TEXT NOT NULL,
			container_id TEXT NOT NULL UNIQUE
		);
		CREATE TABLE IF NOT EXISTS sidecars(
			id TEXT NOT NULL UNIQUE PRIMARY KEY,
			container_id TEXT NOT NULL UNIQUE,
			server_id TEXT NOT NULL,
			FOREIGN KEY (server_id) REFERENCES servers(id)
		);
		INSERT INTO images (id, tag) VALUES ('5b3a4946-e16e-4b14-9e85-cf4ed4fbd017', 'localhost/ayeama/panel/server/minecraft:0.0.1-jre21') ON CONFLICT DO NOTHING;
	`)
	if err != nil {
		panic(err)
	}

	runtime, err := runtime.New(runtime.RuntimeType(config.Config.Runtime))
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	imageRepository := repository.NewImageRepository(db)
	imageService := service.NewImageService(imageRepository)
	imageHandler := handler.NewImageHandler(imageService)
	imageHandler.RegisterHandlers(mux)

	sidecarRepository := repository.NewSidecarRepository(db)

	serverRepository := repository.NewServerRepository(db)
	serverService := service.NewServerService(runtime, serverRepository, imageRepository, sidecarRepository)
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
	slog.Info("starting", slog.String("address", s.server.Addr), slog.String("runtime", config.Config.Runtime))

	var eg errgroup.Group

	api := func() error {
		// return s.server.ListenAndServeTLS("cert.pem", "key.pem")
		return s.server.ListenAndServe()
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
