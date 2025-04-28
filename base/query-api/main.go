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

	"github.com/bradfitz/gomemcache/memcache"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

// Classes
type Response struct {
	Timestamp  string `json:"ts"`
	Query      string `json:"query"`
	Result     string `json:"result"`
	ServerIp   string `json:"serverIp"`
	ServerName string `json:"serverName"`
	Id         int    `json:"id"`
}

type Error struct {
	Message   string `json:"message"`
	Timestamp string `json:"ts"`
	Query     string `json:"query"`
}

type Record struct {
	Id         int    `json:"id"`
	Ip         string `json:"ip"`
	Domain     string `json:"domain"`
	ServerIp   string `json:"serverIp"`
	ServerName string `json:"serverName"`
}

// Variables
var db *sql.DB
var tsFormat = time.RFC3339
var mc *memcache.Client

// Prometheus metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

// Helpers

func connectDB() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB") // e.g. "argodb"

	host := os.Getenv("POSTGRES_HOST") // e.g. "localhost"
	port := os.Getenv("POSTGRES_PORT") // e.g. "5432"

	// FOR NOMAD + CONSUL MESH
	// host := "127.0.0.1"
	// port := "5432"

	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	var err error
	db, _ = sql.Open("postgres", conn)
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to database")
}

func connectCache() {
	host := os.Getenv("MEMCACHED_HOST") // e.g. "localhost"
	port := os.Getenv("MEMCACHED_PORT") // e.g. "5432"

	// FOR NOMAD + CONSUL MESH
	// host := "127.0.0.1"
	// port := "11211"

	conn := fmt.Sprintf("%s:%s", host, port)

	mc = memcache.New(conn)
	err := mc.Ping()
	if err != nil {
		log.Fatal("Failed to connect to cache:", err)
	}
	log.Println("Connected to cache")
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

// Updated writeSuccess to accept an ID parameter.
func writeSuccess(w http.ResponseWriter, id int, serverName string, serverIp string, query string, result string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	ts := time.Now().Format(tsFormat)

	resp := Response{
		Timestamp:  ts,
		Query:      query,
		Result:     result,
		ServerIp:   serverIp,
		ServerName: serverName,
		Id:         id,
	}

	json.NewEncoder(w).Encode(resp)
	log.Printf(`Success: id="%d" query="%s" result="%s" serverName="%s" serverIp="%s"`,
		id, query, result, serverName, serverIp)
}

func writeError(w http.ResponseWriter, statusCode int, message string, query string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	ts := time.Now().Format(tsFormat)

	json.NewEncoder(w).Encode(Error{Message: message, Timestamp: ts, Query: query})
	log.Printf(`Error: status="%d" message="%s" query="%s"`, statusCode, message, query)
}

func checkCache(w http.ResponseWriter, req string, key string) bool {
	cache, err := mc.Get(req + ":" + key)
	if err == nil {
		hit := new(Record)
		json.Unmarshal(cache.Value, &hit)
		log.Printf("Cache hit for %s key: %s", req, key)
		writeSuccess(w, hit.Id, hit.ServerName, hit.ServerIp, hit.Domain, hit.Ip)
		return true
	}
	return false
}

func setCache(key string, record Record) {
	json, _ := json.Marshal(record)
	mc.Set(&memcache.Item{Key: key, Value: json, Expiration: 60})
}

// Handlers

func resolveIP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	query := r.URL.Query().Get("query")

	if query == "" || !isValidDNSName(query) {
		writeError(w, http.StatusBadRequest, "Invalid or missing 'query' parameter", query)
		httpRequestsTotal.WithLabelValues("/api/domain2ip", r.Method, "400").Inc()
		httpRequestDuration.WithLabelValues("/api/domain2ip", r.Method).Observe(time.Since(start).Seconds())
		return
	}

	// Cache check
	if checkCache(w, "domain2ip", query) {
		httpRequestsTotal.WithLabelValues("/api/domain2ip", r.Method, "200").Inc()
		httpRequestDuration.WithLabelValues("/api/domain2ip", r.Method).Observe(time.Since(start).Seconds())
		return
	}

	// Retrieve ID as well
	var id int
	var ip, domain, serverName, serverIp string
	err := db.QueryRow(
		`SELECT r.id, r.ip, r.domain, s.name, s.ip
		FROM records r 
		JOIN servers s ON r.server = s.id 
		WHERE r.domain = $1`,
		query,
	).Scan(&id, &ip, &domain, &serverName, &serverIp)

	if err != nil {
		writeError(w, http.StatusNotFound, "No IP addresses found", query)
		httpRequestsTotal.WithLabelValues("/api/domain2ip", r.Method, "404").Inc()
		httpRequestDuration.WithLabelValues("/api/domain2ip", r.Method).Observe(time.Since(start).Seconds())
		return
	}

	// Cache set
	record := Record{Id: id, Ip: ip, Domain: domain, ServerIp: serverIp, ServerName: serverName}
	setCache("domain2ip:"+query, record)

	// Notice we pass record.Id, record.ServerName, record.ServerIp, ...
	writeSuccess(w, record.Id, serverName, serverIp, domain, ip)
	httpRequestsTotal.WithLabelValues("/api/domain2ip", r.Method, "200").Inc()
	httpRequestDuration.WithLabelValues("/api/domain2ip", r.Method).Observe(time.Since(start).Seconds())
}

func resolveDomain(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	query := r.URL.Query().Get("query")

	if query == "" || !isValidIP(query) {
		writeError(w, http.StatusBadRequest, "Invalid or missing 'query' parameter", query)
		httpRequestsTotal.WithLabelValues("/api/ip2domain", r.Method, "400").Inc()
		httpRequestDuration.WithLabelValues("/api/ip2domain", r.Method).Observe(time.Since(start).Seconds())
		return
	}

	// Cache check
	if checkCache(w, "ip2domain", query) {
		httpRequestsTotal.WithLabelValues("/api/ip2domain", r.Method, "200").Inc()
		httpRequestDuration.WithLabelValues("/api/ip2domain", r.Method).Observe(time.Since(start).Seconds())
		return
	}

	var id int
	var ip, domain, serverName, serverIp string
	err := db.QueryRow(
		`SELECT r.id, r.ip, r.domain, s.name, s.ip
		FROM records r 
		JOIN servers s ON r.server = s.id 
		WHERE r.ip = $1`,
		query,
	).Scan(&id, &ip, &domain, &serverName, &serverIp)

	if err != nil {
		writeError(w, http.StatusNotFound, "No domains found", query)
		httpRequestsTotal.WithLabelValues("/api/ip2domain", r.Method, "404").Inc()
		httpRequestDuration.WithLabelValues("/api/ip2domain", r.Method).Observe(time.Since(start).Seconds())
		return
	}

	// Cache set
	record := Record{Id: id, Ip: ip, Domain: domain, ServerIp: serverIp, ServerName: serverName}
	setCache("ip2domain:"+query, record)

	writeSuccess(w, record.Id, serverName, serverIp, ip, domain)
	httpRequestsTotal.WithLabelValues("/api/ip2domain", r.Method, "200").Inc()
	httpRequestDuration.WithLabelValues("/api/ip2domain", r.Method).Observe(time.Since(start).Seconds())
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
	connectCache()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/ip2domain", resolveDomain)
	mux.HandleFunc("/api/domain2ip", resolveIP)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", getHealth)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}).Handler(mux)

	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))
}
