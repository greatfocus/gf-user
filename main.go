package main

import (
	"os"

	frame "github.com/greatfocus/gf-frame"
	"github.com/greatfocus/gf-user/router"
	"github.com/greatfocus/gf-user/task"
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
	tasks := task.Tasks{}
	tasks.Init(&server)
	server.Cron.Every(1).Sunday().At("19:00:00").Do(tasks.RunDatabaseScripts)
	server.Cron.Every(1).Monday().At("01:00:00").Do(tasks.RebuildIndexes)
	server.Cron.Start()

	// start API service
	server.Start(router.Router(&server))
}
