package platform

import (
	"encoding/json"
	"net/http"
	"telemetry-agent/internal/tracker"
)

func getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracker.GetAll())
}

func NewCollectorApi() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/data", getData)
	return mux
}
