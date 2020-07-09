package main

import (
	frame "github.com/greatfocus/gf-frame"
	"github.com/greatfocus/gf-user/router"
	_ "github.com/lib/pq"
)

// Entry point to the solution
func main() {
	// Load configurations
	service := frame.Create("dev.json")

	// start API service
	service.Start(router.Router(service.DB))
}
