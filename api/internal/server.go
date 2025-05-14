package internal

import (
	"log/slog"
	"net/http"

	"github.com/ayeama/panel/api/internal/broker"
	"github.com/ayeama/panel/api/internal/database"
	"github.com/ayeama/panel/api/internal/handler"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/ayeama/panel/api/internal/service"
)

type Server struct {
	server http.Server
}

func NewServer() *Server {
	database := database.NewDatabase()

	broker, err := broker.NewBroker(broker.BrokerTypeRedis)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	serverRepository := repository.NewServerRepository(database.Db) // todo .Db ugh
	serverService := service.NewServerService(broker, serverRepository)
	serverHandler := handler.NewServerHandler(serverService)
	serverHandler.RegisterHandlers(mux)

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
