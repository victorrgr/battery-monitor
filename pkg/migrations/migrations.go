package migrations

import (
	"database/sql"
	"log"
)

func Run(db *sql.DB) {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS battery_log (
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
            percent INTEGER,
            status TEXT
        )
    `)
	if err != nil {
		log.Fatal("Wasn't able to create migrations for table \"battery_log\"", err)
	}
	log.Println("Migrations ran successfully")
}
