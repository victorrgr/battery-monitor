package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/victorrgr/battery-monitor/pkg/migrations"
	"github.com/victorrgr/battery-monitor/pkg/monitor"
	"log"
)

func closeDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal("Wasn't able to properly close database", err)
	}
}

func main() {
	log.Println("Init")
	db, err := sql.Open("sqlite3", "./battery-monitor.db")
	if err != nil {
		log.Fatal("Error on opening connection to the database", err)
	}
	defer closeDatabase(db)

	migrations.Run(db)
	monitor.Start(db)
}
