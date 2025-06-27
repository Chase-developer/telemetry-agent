package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"telemetry-agent/internal/collector"
	"telemetry-agent/internal/config"
	"telemetry-agent/internal/tracker"
)

func newReverseProxy(targetHost string, targetPort string) *httputil.ReverseProxy {
	target, _ := url.Parse("http://" + targetHost + ":" + targetPort)

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Intercept and modify outgoing request
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req) // preserve default behavior

		tracker.Request(req.URL.Path, req.Method)
		// 🛠 Modify headers or log
		log.Printf("Proxying %s request to %s", req.Method, req.URL.String())

		// Example: Add custom header
		//req.Header.Set("X-Proxy-Intercept", "true")
	}

	// Intercept and inspect response
	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Printf("Got response with status: %d from %s", resp.StatusCode, resp.Request.URL)
		tracker.Response(resp.Request.URL.Path, resp.StatusCode)
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

func Start(cfg *config.Config) {
	authenticate()
	proxy := newReverseProxy(cfg.Backend.ForwardHost, cfg.Backend.ForwardPort)
	address := fmt.Sprintf("%s:%s", cfg.Backend.ListenHost, cfg.Backend.ListenPort)

	go func() {
		log.Printf("Agent proxy listening on %s, forwarding to %s:%s", address, cfg.Backend.ForwardHost, cfg.Backend.ForwardPort)
		err := http.ListenAndServe(address, proxy)
		if err != nil {
			log.Fatalf("Agent ListenAndServe failed: %v", err)
		}
	}()
}

func authenticate() {
	data := map[string]string{
		"id":          "somthing",
		"private_key": "key",
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := http.Post(
		"http://"+collector.Address+"/auth",
		"application/json",
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		log.Println("POST request failed:", err)
		return
	}
	defer resp.Body.Close()

	//log.Println("Response status:", resp.Status, resp.Body)

	var raw map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil {
		log.Println("Failed to decode response:", err)
		return
	}
	log.Printf("Full Response: %+v\n", raw)
}
