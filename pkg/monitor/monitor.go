package monitor

import (
	"database/sql"
)

var running = false

// Start will monitor the battery details and save to the database
func Start(db *sql.DB) {
	running = true
	for running {
		// TODO: Gather data and save to the database
	}
}

// Stop will stop the running monitor
func Stop() {
	running = false
}
