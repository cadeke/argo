package main

// Imports
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

// Classes
type Error struct {
	Timestamp string `json:"ts"`
	Message   string `json:"message"`
}

type Record struct {
	Id     int    `json:"id"`
	Ip     string `json:"ip"`
	Domain string `json:"domain"`
	Server string `json:"server"`
}

// Variables
var db *sql.DB
var tsFormat = time.RFC3339

// Helpers

func connectDB() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST") // e.g. "localhost"
	port := os.Getenv("POSTGRES_PORT") // e.g. "5432"
	dbName := os.Getenv("POSTGRES_DB") // e.g. "argodb"

	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	var err error
	db, _ = sql.Open("postgres", conn)
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to database")
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func isValidDNSName(name string) bool {
	if len(name) == 0 {
		return false
	}

	// Split the domain into labels
	labels := strings.Split(name, ".")

	// Must have at least two labels (name and TLD)
	if len(labels) < 2 {
		return false
	}

	// Check each label
	for _, label := range labels {
		// Check length constraints
		if len(label) == 0 || len(label) > 63 {
			return false
		}

		// Check for hyphens at start or end
		if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return false
		}

		// Check for valid characters
		if !regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(label) {
			return false
		}
	}

	return true
}

func writeSuccess(w http.ResponseWriter, action string, statusCode int, record Record) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(record)
	log.Printf(`%s: id="%d" server="%s" domain="%s" ip="%s"`, action, record.Id, record.Server, record.Domain, record.Ip)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	ts := time.Now().Format(tsFormat)

	json.NewEncoder(w).Encode(Error{Message: message, Timestamp: ts})
	log.Printf(`Error: status="%d" message="%s"`, statusCode, message)
}

// Handlers

func addRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	var record Record
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	if !isValidIP(record.Ip) {
		writeError(w, http.StatusBadRequest, "Invalid IP address: "+record.Ip)
		return
	}

	if !isValidDNSName(record.Domain) {
		writeError(w, http.StatusBadRequest, "Invalid domain name: "+record.Domain)
		return
	}

	var serverId int
	err = db.QueryRow(`SELECT id FROM servers WHERE ip = $1`, record.Server).Scan(&serverId)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Server does not exist: "+record.Server)
		return
	}

	var exists bool
	err = db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM records WHERE ip = $1 AND domain = $2 AND server = $3)`,
		record.Ip, record.Domain, serverId,
	).Scan(&exists)

	if exists {
		writeError(w, http.StatusBadRequest, "Record already exists "+fmt.Sprintf("IP: %s, Domain: %s", record.Ip, record.Domain))
		return
	}

	_, err = db.Exec(
		`INSERT INTO records (ip, domain, server) VALUES ($1, $2, $3)`,
		record.Ip, record.Domain, serverId,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to insert record: "+err.Error())
		return
	}

	var newRecord Record

	db.QueryRow(
		`SELECT id, ip, domain, server FROM records WHERE ip = $1 AND domain = $2 AND server = $3`,
		record.Ip, record.Domain, serverId,
	).Scan(&newRecord.Id, &newRecord.Ip, &newRecord.Domain, &newRecord.Server)

	writeSuccess(w, "ADD", http.StatusCreated, newRecord)
}

func updateRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "Only PUT method is allowed")
		return
	}

	var record Record

	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	// Validate IP and Domain
	if !isValidIP(record.Ip) {
		writeError(w, http.StatusBadRequest, "Invalid IP address: "+record.Ip)
		return
	}
	if !isValidDNSName(record.Domain) {
		writeError(w, http.StatusBadRequest, "Invalid domain name: "+record.Domain)
		return
	}

	// Confirm the server exists
	var serverId int
	err = db.QueryRow(`SELECT id FROM servers WHERE ip = $1`, record.Server).Scan(&serverId)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Server does not exist: "+record.Server)
		return
	}

	// Check if the record with the given ID exists
	var existingId int
	err = db.QueryRow(`SELECT id FROM records WHERE id = $1`, record.Id).Scan(&existingId)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "Record not found "+fmt.Sprintf("ID: %d", record.Id))
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "Error retrieving record: "+err.Error())
		return
	}

	// Update the record’s IP, domain, and server
	_, err = db.Exec(
		`UPDATE records SET ip = $1, domain = $2, server = $3 WHERE id = $4`,
		record.Ip, record.Domain, serverId, record.Id,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update record: "+err.Error())
		return
	}

	// Return the updated record
	var updatedRecord Record
	err = db.QueryRow(
		`SELECT r.id, r.ip, r.domain, s.ip as server 
		 FROM records r 
		 JOIN servers s ON r.server = s.id 
		 WHERE r.id = $1`,
		record.Id,
	).Scan(&updatedRecord.Id, &updatedRecord.Ip, &updatedRecord.Domain, &updatedRecord.Server)

	writeSuccess(w, "UPDATE", http.StatusOK, updatedRecord)
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Only DELETE method is allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing record ID parameter")
		return
	}

	// Check if the record with the given ID exists
	var existingId int
	err := db.QueryRow(`SELECT id FROM records WHERE id = $1`, id).Scan(&existingId)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "Record not found "+fmt.Sprintf("ID: %d", existingId))
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "Error retrieving record: "+err.Error())
		return
	}

	// Delete the record
	_, err = db.Exec(`DELETE FROM records WHERE id = $1`, existingId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete record: "+err.Error())
		return
	}

	writeSuccess(w, "DELETE", http.StatusOK, Record{Id: existingId})
}

func getAllRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
		return
	}

	rows, err := db.Query(`SELECT r.id, r.ip, r.domain, s.ip as server 
		FROM records r 
		JOIN servers s ON r.server = s.id`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve records: "+err.Error())
		return
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var record Record
		if err := rows.Scan(&record.Id, &record.Ip, &record.Domain, &record.Server); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to scan record: "+err.Error())
			return
		}
		records = append(records, record)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(records)
}

func getRecordById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Missing record ID parameter")
		return
	}

	var record Record
	err := db.QueryRow(
		`SELECT r.id, r.ip, r.domain, s.ip as server 
		 FROM records r 
		 JOIN servers s ON r.server = s.id 
		 WHERE r.id = $1`, id,
	).Scan(&record.Id, &record.Ip, &record.Domain, &record.Server)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "Record not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve record: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(record)
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Create proper JSON response structure
	response := struct {
		Status string `json:"status"`
	}{
		Status: "Healthy!",
	}

	// Handle encoding error
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding health response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Main

func main() {

	connectDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/add", addRecord)
	mux.HandleFunc("/api/update", updateRecord)
	mux.HandleFunc("/api/delete", deleteRecord)
	mux.HandleFunc("/api/get", getRecordById)
	mux.HandleFunc("/api/list", getAllRecords)
	mux.HandleFunc("/health", getHealth)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}).Handler(mux)

	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
