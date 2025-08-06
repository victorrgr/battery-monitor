package analyser

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/victorrgr/battery-monitor/pkg/monitor"
	"github.com/victorrgr/battery-monitor/pkg/system"
	"github.com/victorrgr/battery-monitor/pkg/utils"
	"github.com/victorrgr/battery-monitor/templates"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Data struct {
	List []monitor.BatteryLog
}

func Analyze(db *sql.DB) {
	http.Handle("/", http.FileServer(http.FS(templates.Files)))
	http.HandleFunc("/dates", datesHandler(db))
	http.HandleFunc("/data", dataHandler(db))
	log.Println("Opening Server at 8080")
	exec.Command("xdg-open", "http://localhost:8080/report.html")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Unable to open server at port 8080", err)
	}

	//list := searchData(db)
	//list = sample(list, 100)
	//generateReport(Data{List: list})
	//log.Println("Generated report.html")
}

func dataHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var date = time.Now()
		if r.URL.Query().Has("date") {
			date, err = time.Parse("2006-01-02", r.URL.Query().Get("date"))
			if err != nil {
				sendBadRequest(w, `Invalid date format for query param "date": `+err.Error())
				return
			}
		}

		list := searchData(db, date)
		list = sample(list, 70)
		marshal, err := json.Marshal(list)
		if err != nil {
			msg := "Error Transforming to JSON: " + err.Error()
			_, _ = w.Write([]byte(msg))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(marshal)
	}
}

func sendInternalServerError(w http.ResponseWriter, msg string) {
	sendError(w, msg, http.StatusInternalServerError)
}

func sendBadRequest(w http.ResponseWriter, msg string) {
	sendError(w, msg, http.StatusBadRequest)
}

func sendError(w http.ResponseWriter, msg string, statusCode int) {
	data := map[string]string{
		"message": msg,
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		msg := "Error Transforming to JSON: " + err.Error()
		_, _ = w.Write([]byte(msg))
	}
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(marshal)
}

func datesHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var page int32 = 0
		query := r.URL.Query()
		if query.Has("page") {
			page, err = utils.ParseInt32(query.Get("page"))
			if err != nil {
				sendBadRequest(w, `Error Parsing Request Param "page": `+err.Error())
				return
			}
		}
		var size int32 = 5
		if query.Has("size") {
			size, err = utils.ParseInt32(query.Get("size"))
			if err != nil {
				sendBadRequest(w, `Error Parsing Request Param "page": `+err.Error())
				return
			}
		}

		offset := (page - 1) * size
		datesRes, err := searchDates(db, size, offset)
		if err != nil {
			sendInternalServerError(w, "Error Fetching Data: "+err.Error())
		}

		marshal, err := json.Marshal(datesRes)
		if err != nil {
			msg := "Error Transforming to JSON: " + err.Error()
			_, _ = w.Write([]byte(msg))
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(marshal)
	}
}

type DatesResponse struct {
	Dates      []string `json:"dates"`
	TotalItems int      `json:"totalItems"`
	TotalPages int      `json:"totalPages"`
}

func searchDates(db *sql.DB, size int32, offset int32) (DatesResponse, error) {
	var response DatesResponse
	countQuery := `
		SELECT COUNT(*) FROM (
			SELECT DATE(DATETIME("timestamp", 'localtime')) AS day
			FROM battery_log
			GROUP BY day
		)
	`
	err := db.QueryRow(countQuery).Scan(&response.TotalItems)
	if err != nil {
		return response, fmt.Errorf("error counting total items: %w", err)
	}

	if size > 0 {
		response.TotalPages = (response.TotalItems + int(size) - 1) / int(size)
	} else {
	}

	// 3. Fetch paginated dates
	query := `
		SELECT DATE(DATETIME("timestamp", 'localtime')) AS day
		FROM battery_log
		GROUP BY day
		ORDER BY day
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, size, offset)
	if err != nil {
		return response, fmt.Errorf("error querying paginated dates: %w", err)
	}

	for rows.Next() {
		var day string
		if err := rows.Scan(&day); err != nil {
			return response, fmt.Errorf("error scanning date: %w", err)
		}
		response.Dates = append(response.Dates, day)
	}

	return response, nil
}

func sample(list []monitor.BatteryLog, maxPoints int) []monitor.BatteryLog {
	if len(list) > maxPoints {
		step := len(list) / maxPoints
		var sampled []monitor.BatteryLog
		for i := 0; i < len(list); i += step {
			sampled = append(sampled, list[i])
		}
		list = sampled
	}
	return list
}

func generateReport(data any) {
	tmpl, err := template.ParseFS(templates.Files, "report.html")
	if err != nil {
		log.Fatal("Template parsing error: ", err)
	}

	file, err := os.Create("report.html")
	if err != nil {
		log.Fatal("Could not create output file:", err)
	}
	defer system.CloseFile(file)

	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatal("Error executing template: ", err)
	}
}

func searchData(db *sql.DB, date time.Time) []monitor.BatteryLog {
	query := `
	SELECT "timestamp", percent, status FROM battery_log
	WHERE DATE(DATETIME("timestamp", 'localtime')) = DATE(?)
	ORDER BY "timestamp";
	`
	res, err := db.Query(query, date.Format("2006-01-02"))
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
