package analyser

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/victorrgr/battery-monitor/pkg/monitor"
	"github.com/victorrgr/battery-monitor/pkg/system"
	"github.com/victorrgr/battery-monitor/templates"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

type Data struct {
	List []monitor.BatteryLog
}

func Analyze(db *sql.DB, port int) {
	http.Handle("/", http.FileServer(http.FS(templates.Files)))
	http.HandleFunc("/dates", datesHandler(db))
	http.HandleFunc("/data", dataHandler(db))
	host := fmt.Sprintf("localhost:%d", port)
	log.Printf("Web server starting at http://%s\n", host)
	exec.Command("xdg-open", fmt.Sprintf("http://%s", host))
	err := http.ListenAndServe(host, nil)
	if err != nil {
		if errors.Is(err, syscall.EADDRINUSE) {
			log.Fatalf("Port %d already in use", port)
		}
		log.Fatalf("Unable to open web server at port %d: %s", port, err)
	}
}

func dataHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var date = time.Now()
		query := r.URL.Query()
		if query.Has("date") {
			date, err = time.Parse("2006-01-02", query.Get("date"))
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
			_, err = w.Write([]byte(msg))
			if err != nil {
				log.Println("[ERROR] Writing error response: ", err)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(marshal)
		if err != nil {
			log.Println("[ERROR] Writing error response: ", err)
		}
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
		var page = 0
		query := r.URL.Query()
		if query.Has("page") {
			page, err = strconv.Atoi(query.Get("page"))
			if err != nil {
				sendBadRequest(w, `Error Parsing Request Param "page": `+err.Error())
				return
			}
			if page < 0 {
				sendBadRequest(w, `Request Param "page" cannot be less than 0`)
				return
			}
		}
		var size = 5
		if query.Has("size") {
			size, err = strconv.Atoi(query.Get("size"))
			if err != nil {
				sendBadRequest(w, `Error Parsing Request Param "page": `+err.Error())
				return
			}
			if size < 0 {
				sendBadRequest(w, `Request Param "size" cannot be less than 0`)
				return
			}
			if size == 0 {
				sendBadRequest(w, `Request Param "size" cannot be equals to 0`)
				return
			}
		}

		offset := page * size
		datesRes, err := searchDates(db, size, offset)
		if err != nil {
			sendInternalServerError(w, "Error Fetching Data: "+err.Error())
		}

		if datesRes.TotalPages == page {
			sendError(w, "Requested page exceeds total available pages", http.StatusUnprocessableEntity)
			return
		}

		marshal, err := json.Marshal(datesRes)
		if err != nil {
			msg := "Error Transforming to JSON: " + err.Error()
			_, err = w.Write([]byte(msg))
			if err != nil {
				log.Println("[ERROR] Writing error response: ", err)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(marshal)
		if err != nil {
			log.Println("[ERROR] Writing error response: ", err)
			return
		}
	}
}

type DatesResponse struct {
	TotalPages int      `json:"totalPages"`
	TotalItems int      `json:"totalItems"`
	Dates      []string `json:"dates"`
}

func searchDates(db *sql.DB, size int, offset int) (DatesResponse, error) {
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
		response.TotalPages = (response.TotalItems + size - 1) / size
	}

	query := `
		SELECT DATE(DATETIME("timestamp", 'localtime')) AS day
		FROM battery_log
		GROUP BY day
		ORDER BY day DESC
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

	if response.Dates == nil {
		response.Dates = make([]string, 0)
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
	tmpl, err := template.ParseFS(templates.Files, "index.html")
	if err != nil {
		log.Fatal("Template parsing error: ", err)
	}

	file, err := os.Create("index.html")
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

func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
