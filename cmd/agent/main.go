package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Create a reverse proxy to http://localhost:8081
func newReverseProxy(targetPort string) *httputil.ReverseProxy {
	target, _ := url.Parse("http://localhost" + targetPort)

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Intercept and modify outgoing request
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req) // preserve default behavior

		// 🛠 Modify headers or log
		log.Printf("Proxying %s request to %s", req.Method, req.URL.String())

		// Example: Add custom header
		//req.Header.Set("X-Proxy-Intercept", "true")
	}

	// Intercept and inspect response
	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Printf("Got response with status: %d from %s", resp.StatusCode, resp.Request.URL)

		// Example: Inject response header
		//resp.Header.Set("X-Proxy-Processed", "yes")

		// Example: Log body length
		log.Printf("Response body size: %d", resp.ContentLength)

		// Optionally modify body (advanced — involves re-reading it)
		return nil
	}

	// Optional: Custom error handling
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		http.Error(w, "Proxy error: "+err.Error(), http.StatusBadGateway)
	}

	return proxy
}

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

func main() {
	proxy := newReverseProxy(":8081")

	// log.Println("Listening on :8080 and proxying to :8081")
	// err := http.ListenAndServe(":8080", proxy)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	go func() {
		log.Println("Listening on :8080 and proxying to :8081")
		err := http.ListenAndServe(":8080", proxy)
		if err != nil {
			log.Fatal("Proxy error: ", err)
		}
	}()
	// Create separate routers
	// mainMux := http.NewServeMux()
	// mainMux.HandleFunc("/", rootHandler)
	// mainMux.HandleFunc("/ping", pingHandler)
	// mainMux.HandleFunc("/hello", helloHandler)

	puppetMux := http.NewServeMux()
	puppetMux.HandleFunc("/cmd", puppetCommandHandler)
	puppetMux.HandleFunc("/status", puppetStatusHandler)
	puppetMux.HandleFunc("/", rootHandler)
	puppetMux.HandleFunc("/ping", pingHandler)
	puppetMux.HandleFunc("/hello", helloHandler)

	// Start main server (port 8080)

	// go func() {
	// 	log.Println("Main app listening on http://localhost:8080")
	// 	err := http.ListenAndServe(":8080", mainMux)
	// 	if err != nil {
	// 		log.Fatal("Main app error: ", err)
	// 	}
	// }()

	// Start puppet API server (port 8081)
	go func() {
		log.Println("Puppet API listening on http://localhost:8081")
		err := http.ListenAndServe(":8081", puppetMux)
		if err != nil {
			log.Fatal("Puppet API error: ", err)
		}
	}()

	// Block main goroutine
	select {}
}
