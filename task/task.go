package task

import (
	"log"

	"github.com/greatfocus/gf-frame/server"

	"github.com/greatfocus/gf-frame/config"
	"github.com/greatfocus/gf-frame/database"
)

// Tasks struct
type Tasks struct {
	config *config.Config
	db     *database.Conn
}

// Init required parameters
func (t *Tasks) Init(server *server.Server) {
	t.config = server.Config
	t.db = server.DB
}

// RunDatabaseScripts intiates running database scripts
func (t *Tasks) RunDatabaseScripts() {
	log.Println("Scheduler_RunDatabaseScripts started")
	if t.config.Database.Master.ExecuteSchema {
		t.db.Master.ExecuteSchema(t.db.Master.Conn)
	}
	if t.config.Database.Slave.ExecuteSchema {
		t.db.Slave.ExecuteSchema(t.db.Slave.Conn)
	}
	log.Println("Scheduler_RunDatabaseScripts ended")
}

// RebuildIndexes make changes to indexes
func (t *Tasks) RebuildIndexes() {
	log.Println("Scheduler_RebuildIndexes started")
	if t.config.Database.Master.ExecuteSchema {
		t.db.Master.RebuildIndexes(t.db.Master.Conn, "gf_user")
	}
	if t.config.Database.Slave.ExecuteSchema {
		t.db.Slave.RebuildIndexes(t.db.Slave.Conn, "gf_user")
	}
	log.Println("Scheduler_RebuildIndexes ended")
}
