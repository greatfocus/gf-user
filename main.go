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
	service := frame.Create("dev.json")

	// background task
	task := task.Task{}
	task.Init(service.DB, service.Config)
	cron.Every(10).Second().Do(task.SendNotification)
	cron.Start()

	// start API service
	service.Start(router.Router(service.DB))
}
