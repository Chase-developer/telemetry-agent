package puppetapi

import (
	"fmt"
	"log"
	"net/http"
	"telemetry-agent/internal/config"
)

// --- Main App Handlers (Port 8080) ---
func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Main: ", r.Method, r.URL.Path)
	fmt.Fprintln(w, "Welcome to the Go server!")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "stranger"
	}
	fmt.Fprintf(w, "Hello, %s!\n", name)
}

// --- Puppet API Handlers (Port 8081) ---
func puppetCommandHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("PuppetAPI: ", r.Method, r.URL.Path)
	fmt.Fprintln(w, "Received puppet command")
}

func puppetStatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Puppet API status: OK")
}

// --- Exported Router Constructor ---
func newPuppetAPI() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/cmd", puppetCommandHandler)
	mux.HandleFunc("/status", puppetStatusHandler)
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/hello", helloHandler)
	return mux
}

func Start(cfg *config.Config) {
	puppetApi := newPuppetAPI()
	address := fmt.Sprintf("%s:%s", cfg.Backend.ForwardHost, cfg.Backend.ForwardPort)
	// Start puppet API server (port 8081)
	go func() {
		log.Printf("Puppet API listening on %s", address)
		err := http.ListenAndServe(address, puppetApi)
		if err != nil {
			log.Fatal("Puppet API error: ", err)
		}
	}()
}
