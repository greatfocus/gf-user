package main

import (
	frame "github.com/greatfocus/gf-frame"
	"github.com/greatfocus/gf-user/router"
	_ "github.com/lib/pq"
)

// Entry point to the solution
func main() {
	// Load configurations
	server := frame.Create("dev.json")
	server.Cron.Start()

	// start API service
	server.Start(router.Router(&server))
}
