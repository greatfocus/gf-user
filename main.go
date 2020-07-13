package main

import (
	frame "github.com/greatfocus/gf-frame"
	"github.com/greatfocus/gf-frame/cron"
	"github.com/greatfocus/gf-user/router"
	"github.com/greatfocus/gf-user/task"
	_ "github.com/lib/pq"
)

// Entry point to the solution
func main() {
	// Load configurations
	server := frame.Create("dev.json")

	// background task
	task := task.Task{}
	task.Init(server.DB, server.Config)
	cron.Every(30).Second().Do(task.SendNotification)
	cron.Start()

	// start API service
	server.Start(router.Router(server.DB, server.Config))
}
