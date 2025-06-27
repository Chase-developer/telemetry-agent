package collector

import (
	"encoding/json"
	"log"
	"net/http"
	"telemetry-agent/internal/config"
	"telemetry-agent/internal/tracker"
)

const Address = "localhost:8082"

func getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracker.GetAll())
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  200,
		"message": "Authenticate",
	}
	json.NewEncoder(w).Encode(response)
}

func newCollectorApi() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/data", getData)
	mux.HandleFunc("/auth", authenticate)
	return mux
}

func Start(cfg *config.Config) {
	collectorApi := newCollectorApi()
	// Start puppet API server (port 8081)
	go func() {
		log.Printf("Collector API listening on %s", Address)
		err := http.ListenAndServe(Address, collectorApi)
		if err != nil {
			log.Fatal("Collector API error: ", err)
		}
	}()
}
