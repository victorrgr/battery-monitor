package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/victorrgr/battery-monitor/pkg/system"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/victorrgr/battery-monitor/pkg/analyser"
	"github.com/victorrgr/battery-monitor/pkg/migrations"
	"github.com/victorrgr/battery-monitor/pkg/monitor"
)

func getDatabasePath() (string, error) {
	dataDir := system.GetSharedLocalDir()
	dbDir := filepath.Join(dataDir, "battery-monitor")
	err := os.MkdirAll(dbDir, 0755)
	if err != nil {
		return "", err
	}

	return filepath.Join(dbDir, "battery-monitor.db"), nil
}

func closeDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal("Wasn't able to properly close database connection", err)
	}
}

func usageStr() string {
	return `Usage: battery-monitor <command> [options]

Commands:
  monitor                Start battery monitoring
  analyze --port=<port>  Analyze data and start web server on given port
  migrate                Run DB migrations`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usageStr())
		return
	}
	dataSourceName, err := getDatabasePath()
	if err != nil {
		log.Fatal("Error determining database path location: ", err)
	}
	db, err := sql.Open("sqlite3", dataSourceName)
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
		flagSet := flag.NewFlagSet(cmd, flag.ExitOnError)
		port := flagSet.Int("port", 8080, "Port to host analyser server")
		err := flagSet.Parse(os.Args[2:])
		if err != nil {
			log.Fatal("Error parsing arguments: ", err)
		}
		analyser.Analyze(db, *port)
	case "migrate":
		migrations.Run(db)
	default:
		fmt.Println("unknown command:", cmd)
		fmt.Println(usageStr())
	}
}
