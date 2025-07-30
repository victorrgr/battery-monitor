package analyser

import (
	"database/sql"
	"log"
)

func Analyze(db *sql.DB) {
	log.Fatal("Not implemented")

	//res, err := db.Query("SELECT timestamp, percent, status FROM battery_log")
	//if err != nil {
	//	panic(err)
	//}
	//
	//log.Println("queried table")
	//
	//for res.Next() {
	//	var timestamp *time.Time
	//	var percent int32
	//	var status string
	//	err := res.Scan(&timestamp, &percent, &status)
	//	if err != nil {
	//		panic(err)
	//	}
	//	log.Printf("timestamp: %s", timestamp)
	//	log.Printf("percent: %d", percent)
	//	log.Printf("status: %s", status)
	//}
}
