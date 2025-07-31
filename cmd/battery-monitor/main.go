package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/victorrgr/battery-monitor/pkg/analyser"
	"github.com/victorrgr/battery-monitor/pkg/migrations"
	"github.com/victorrgr/battery-monitor/pkg/monitor"
)

func closeDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal("Wasn't able to properly close database connection", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: battery-monitor [monitor|analyze|migrate]")
		return
	}
	db, err := sql.Open("sqlite3", "./battery-monitor.db")
	if err != nil {
		log.Fatal("Error on opening connection to the database", err)
	}
	defer closeDatabase(db)

	cmd := os.Args[1]

	switch cmd {
	case "monitor":
		migrations.Run(db)
		monitor.Start(db)
	case "analyze", "analyse":
		analyser.Analyze(db)
	case "migrate":
		migrations.Run(db)
	default:
		fmt.Println("unknown command:", cmd)
		fmt.Println("usage: battery-monitor [monitor|analyze|migrate]")
	}
}
