package main

import (
	"os"

	frame "github.com/greatfocus/gf-frame"
	"github.com/greatfocus/gf-user/router"
	_ "github.com/lib/pq"
)

// Entry point to the solution
func main() {
	// Get arguments
	if os.Args[1] == "" {
		panic("Pass the environment")
	}

	// Load configurations
	server := frame.Create(os.Args[1] + ".json")
	server.Cron.Start()

	// start API service
	server.Start(router.Router(&server))
}
