package main

import (
	"os"

	frame "github.com/greatfocus/gf-frame"
	"github.com/greatfocus/gf-user/router"
	_ "github.com/greatfocus/pq"
)

// Entry point to the solution
func main() {
	// Get arguments
	if os.Args[1] == "" {
		panic("Pass the environment")
	}

	// Load configurations
	server := frame.Create(os.Args[1] + ".json")

	// background task
	server.Cron.Start()

	// start API service
	server.Start(router.Router(&server))
}
