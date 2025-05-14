package main

import internal "github.com/ayeama/panel/agent/internal"

func main() {
	// slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	agent := internal.NewAgent()
	agent.Start()
}
