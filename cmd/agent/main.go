package main

import (
	"log"

	"telemetry-agent/internal/agent"
	"telemetry-agent/internal/collector"
	"telemetry-agent/internal/config"
	"telemetry-agent/internal/puppetapi"
)

/*
* tasks:
create a configuration file to read from which says which port to listen and which port to forward to.
probably don't need the url
record the time stamp of the url and store it in a log file or smth
creating another file which runs the collector api. can probbaly just set it to port 555 or smth.
so that the class can be easily imported. the configuration file will probably also have that setting
wait maybe not, tho I think can be as long as the url of the official deployment data collector doesn't change
*/

// Create a reverse proxy to http://localhost:8081

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	log.Println(cfg)
	log.Println(err)
	collector.Start(cfg)
	agent.Start(cfg)

	puppetapi.Start(cfg)
	// proxy := agent.NewReverseProxy(cfg.Backend.ForwardHost, cfg.Backend.ForwardPort)

	// // log.Println("Listening on :8080 and proxying to :8081")
	// // err := http.ListenAndServe(":8080", proxy)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// go func() {
	// 	log.Println("Listening on ", cfg.Backend.ListenHost+":"+cfg.Backend.ListenPort,
	// 		"and proxying to", cfg.Backend.ForwardHost+":"+cfg.Backend.ForwardPort)
	// 	err := http.ListenAndServe(":"+cfg.Backend.ListenPort, proxy)
	// 	if err != nil {
	// 		log.Fatal("Proxy error: ", err)
	// 	}
	// }()

	// Create separate routers
	// mainMux := http.NewServeMux()
	// mainMux.HandleFunc("/", rootHandler)
	// mainMux.HandleFunc("/ping", pingHandler)
	// mainMux.HandleFunc("/hello", helloHandler)

	// puppetMux := http.NewServeMux()
	// puppetMux.HandleFunc("/cmd", puppetCommandHandler)
	// puppetMux.HandleFunc("/status", puppetStatusHandler)
	// puppetMux.HandleFunc("/", rootHandler)
	// puppetMux.HandleFunc("/ping", pingHandler)
	// puppetMux.HandleFunc("/hello", helloHandler)

	// Start main server (port 8080)

	// go func() {

	// collectorApi := platform.NewCollectorApi()
	// // Start puppet API server (port 8081)
	// go func() {
	// 	log.Println("Collector API listening on http://localhost:8082")
	// 	err := http.ListenAndServe(":8082", collectorApi)
	// 	if err != nil {
	// 		log.Fatal("Collector API error: ", err)
	// 	}
	// }()

	// Block main goroutine
	select {}
}
