package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type CustomHandler struct {
}
type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func main() {
	var err error
	dsn := "mariom:password@tcp(127.0.0.1:3306)/event_management?parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect db")
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Db is unreachable")
	}
	http.HandleFunc("/events", eventsHandler)
	http.HandleFunc("/event/", eventHandler)
	fmt.Println("Server running on port : 8080")
	http.ListenAndServe(":8080", nil)
}
func eventsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listEvent(w, r)
	case http.MethodPost:
		createEvent(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
func createEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Hello world From Create Event!")
	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now
	result, err := db.Exec("Insert into events(title, description, location, start_time, end_time, created_by, created_at, updated_at) values(?,?,?,?,?,?,?,?)", event.Title, event.Description, event.Location, event.StartTime, event.EndTime, event.CreatedBy, event.CreatedAt, event.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	event.ID = int(id)
	resp := map[string]interface{}{
		"message": "Event Created Successfully",
		"data":    event,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
func eventHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/event/") //find id from url trimming /event/
	id, err := strconv.Atoi(idStr)                     // id convert into int
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case http.MethodGet:
		getEventDetails(w, id)
	case http.MethodPut:
		updateEvent(w, id, r)
	case http.MethodDelete:
		deleteEvent(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	fmt.Fprintln(w, "Hello , world From Events Handler!")
}
func listEvent(w http.ResponseWriter, r *http.Request) {
	var events []Event
	rows, err := db.Query("Select * from events")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.EndTime, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		events = append(events, e)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Println("header: ", w.Header().Get("Content-Type"))
	json.NewEncoder(w).Encode(events)
}
func getEventDetails(w http.ResponseWriter, id int) {
	e, err := getEventFromDb(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}
func updateEvent(w http.ResponseWriter, id int, r *http.Request) {
	event, err := getEventFromDb(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var updateEvent Event
	if err := json.NewDecoder(r.Body).Decode(&updateEvent); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	updateEvent.ID = id
	updateEvent.CreatedAt = event.CreatedAt
	updateEvent.UpdatedAt = time.Now()
	_, err = db.Exec(`UPDATE events SET title=?, description=?, location=?, start_time=?,end_time=?, created_by=?, created_at=?, updated_at=? WHERE id=?`,
		updateEvent.Title, updateEvent.Description, updateEvent.Location, updateEvent.StartTime, updateEvent.EndTime, updateEvent.CreatedBy, updateEvent.CreatedAt, updateEvent.UpdatedAt, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{ // map resp key->string, value(data)->interface means any like string, object, int
		"message": "Event Updated Successfully",
		"data":    updateEvent,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

}
func deleteEvent(w http.ResponseWriter, id int) {

	result, err := db.Exec(`DELETE from events  WHERE id=?`, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		http.Error(w, "No Event Found of this given Id", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Event deleted Successfully")

}
func getEventFromDb(id int) (*Event, error) {
	var e Event
	err := db.QueryRow(`Select * from events where id =?`, id).Scan(&e.ID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.EndTime, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// func (c CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Hello , world From Custom Handler!")
// }
// func helloHandler(w http.ResponseWriter, r *http.Request) {

// 	fmt.Fprintln(w, "Hello , world From HTTP Server!")

// }
// func main() {
// 	//h := &CustomHandler{}
// 	mux := http.NewServeMux()
// 	mux.Handle("/user", userHandler)
// 	mux.Handle("/user", eventHandler)
// 	fmt.Println("Server running on port : 8080")
// 	//	http.ListenAndServe(":8080", h)
// 	http.HandleFunc("/", helloHandler)
// 	http.ListenAndServe(":8080", mux)
// }
