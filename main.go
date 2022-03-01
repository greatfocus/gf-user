package main

import (
	_ "github.com/greatfocus/gf-pq"
	frame "github.com/greatfocus/gf-sframe"
	"github.com/greatfocus/gf-user/router"
)

// Entry point to the solution
func main() {

	server := frame.NewFrame("gf-user")
	mux := router.LoadRouter(server.Server)
	server.Start(mux)
}

// err = conn.StartConsumer("test-queue", "test-key", handler, 2)
// conn.Publish("test-key", []byte(`{"message":"test"}`))
