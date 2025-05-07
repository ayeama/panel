package main

import (
	"github.com/ayeama/panel/api/internal"
)

func main() {
	// slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	a := internal.NewServer()
	a.Start()
}
