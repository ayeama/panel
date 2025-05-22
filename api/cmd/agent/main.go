package main

import "github.com/ayeama/panel/api/internal"

func main() {
	agent := internal.NewAgent()
	agent.Start()
}
