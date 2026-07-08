package main

import "github.com/ayeama/panel/internal"

func main() {
	server := internal.NewServer()
	server.Run()
}
