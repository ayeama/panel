package main

import (
	"github.com/ayeama/panel/api/internal"
	"github.com/ayeama/panel/api/internal/config"
)

func main() {
	// slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	config.New()
	server := internal.NewServer()
	server.Start()
}
