package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Print verbose startup information
	log.Println("========== STARTUP DIAGNOSTICS ==========")
	log.Println("Current working directory:", getwd())
	log.Println("Environment variables:")
	for _, env := range os.Environ() {
		log.Println("  ", env)
	}
	
	// Very simple HTTP server
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health check endpoint called!")
		w.WriteHeader(200)
		fmt.Fprint(w, "OK")
	})
	
	// Listen on both the PORT env var and port 3000 as fallback
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	
	// Log listening information
	log.Printf("Starting minimal server on port %s", port)
	log.Printf("Listening on http://0.0.0.0:%s", port)
	
	// Start server with additional error handling
	server := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 3 * time.Second,
	}
	
	log.Fatal(server.ListenAndServe())
}

func getwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return "ERROR: " + err.Error()
	}
	return dir
}