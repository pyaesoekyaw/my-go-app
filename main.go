// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	// "os" // You can remove os import if hardcoding port
)

func main() {
	// Let's explicitly use 8000 here
	port := "8000" // Hardcoding for certainty

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from your Go app! Running on port %s\n", port)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	log.Printf("Server starting on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
