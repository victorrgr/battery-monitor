package analyser

import (
	"database/sql"
	"github.com/victorrgr/battery-monitor/pkg/monitor"
	"github.com/victorrgr/battery-monitor/templates"
	"html/template"
	"log"
	"os"
	"time"
)

type ReportData struct {
	Timestamps []string
	Percents   []float64
}

func Analyze(db *sql.DB) {
	list := searchData(db)

	var data ReportData
	for _, entry := range list {
		data.Timestamps = append(data.Timestamps, entry.Timestamp.Format("15:04:05"))
		data.Percents = append(data.Percents, float64(entry.Percent))
	}

	tmpl, err := template.ParseFS(templates.Files, "report.gohtml")
	if err != nil {
		log.Fatal("Template parsing error: ", err)
	}

	file, err := os.Create("report.html")
	if err != nil {
		log.Fatal("Could not create output file:", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatal("Error executing template: ", err)
	}

	log.Println("Generated report.html")
}

func searchData(db *sql.DB) []monitor.BatteryLog {
	res, err := db.Query("SELECT timestamp, percent, status FROM battery_log")
	if err != nil {
		log.Fatal("Error searching data: ", err)
	}

	var list []monitor.BatteryLog
	for res.Next() {
		var timestamp time.Time
		var percent float32
		var status string
		err := res.Scan(&timestamp, &percent, &status)
		if err != nil {
			log.Fatal("Error scanning fields from response: ", err)
		}
		parseStatus, err := monitor.ParseStatus(status)
		if err != nil {
			log.Fatal("Error on parsing status: ", err)
		}
		list = append(list, monitor.BatteryLog{
			Timestamp: timestamp,
			Percent:   percent,
			Status:    parseStatus,
		})
	}
	return list
}
