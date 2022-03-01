package main

import (
	"github.com/greatfocus/gf-sframe/server"
	"github.com/greatfocus/gf-user/router"
	_ "github.com/lib/pq"
)

// main entry point to the service
func main() {
	service := server.NewServer("gf-user", "user")
	service.Mux = router.LoadRouter(service)
	service.Start()
}
