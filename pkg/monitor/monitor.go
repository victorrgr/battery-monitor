package monitor

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Status string

const (
	CHARGING     Status = "Charging"
	DISCHARGING  Status = "Discharging"
	NOT_CHARGING Status = "Not charging"
)

type BatteryLog struct {
	Timestamp time.Time
	Percent   float32
	Status    Status
}

var running = false

// Start will monitor the battery details and save to the database
func Start(db *sql.DB) {
	log.Println("Monitoring")
	running = true
	for running {
		// TODO: Gather data and save to the database
		energyNow := parseInt32(readBatteryField("energy_now"))
		energyFull := parseInt32(readBatteryField("energy_full"))
		percent := float32(energyNow) / float32(energyFull) * 100

		status, err := ParseStatus(readBatteryField("status"))
		if err != nil {
			log.Fatal("Error parsing battery status")
		}

		_, err = save(db, BatteryLog{
			Percent:   percent,
			Status:    status,
			Timestamp: time.Now(),
		})
		if err != nil {
			log.Fatal("Error saving battery log data to the database", err)
		}
		// TODO: Increase to minutes?
		time.Sleep(time.Second)
	}
}

func ParseStatus(s string) (Status, error) {
	switch s {
	case string(CHARGING):
		return CHARGING, nil
	case string(DISCHARGING):
		return DISCHARGING, nil
	case string(NOT_CHARGING):
		return NOT_CHARGING, nil
	default:
		return "", fmt.Errorf("invalid status: %s", s)
	}
}

func parseInt32(str string) int32 {
	parsed, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		log.Fatalf("Error to parse value \"%s\": %s\n", err)
	}
	return int32(parsed)
}

func readBatteryField(field string) string {
	data, err := os.ReadFile("/sys/class/power_supply/BAT0/" + field)
	if err != nil {
		log.Fatalf("Was not able to gather data for field \"%s\": %s\n", field, err)
	}
	return strings.TrimSpace(string(data))
}

func save(db *sql.DB, log BatteryLog) (sql.Result, error) {
	query := `
	INSERT INTO battery_log(timestamp, percent, status)
	VALUES (:timestamp, :percent, :status)
	`
	return db.Exec(query, log.Timestamp, log.Percent, log.Status)
}

// Stop will stop the running monitor
func Stop() {
	running = false
}
