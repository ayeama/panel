package main

import "github.com/ayeama/panel/internal/api"

func main() {
	server := api.NewServer()
	server.Run()
}
