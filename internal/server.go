package internal

import (
	"context"
	"log"
	"net/http"

	"github.com/ayeama/panel/internal/handler"
	"github.com/ayeama/panel/internal/middleware"
	"github.com/ayeama/panel/internal/runtime"
	"go.podman.io/podman/v6/pkg/bindings"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (s *Server) Run() {
	ctx := context.TODO()

	podman, err := bindings.NewConnection(ctx, "unix:///run/user/1000/podman/podman.sock")
	if err != nil {
		log.Fatal(err)
	}

	runtime := runtime.NewPodmanRuntime(&podman)

	mux := http.NewServeMux()

	instanceHandler := handler.NewInstanceHandler(&runtime)
	instanceHandler.RegisterHandlers(mux)

	imageHandler := handler.NewImageHandler(&runtime)
	imageHandler.RegisterHandlers(mux)

	go webhook(&runtime)

	log.Println("starting")

	err = http.ListenAndServe("0.0.0.0:8000", middleware.Log(middleware.Cors(mux)))
	if err != nil {
		log.Fatal(err)
	}
}
