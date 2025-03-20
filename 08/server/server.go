package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request")
		ctx := r.Context()
		// time.Sleep(2 * time.Second) // Simulate a delay
		// fmt.Fprintln(w, "Hello from the Server!")

		select {
			case <-time.After(2 * time.Second):
				fmt.Fprintln(w, "Hello from the Server!")
			case <-ctx.Done():
				err := ctx.Err()
				log.Println("Request cancelled:", err)
				http.Error(w, "Request cancelled", http.StatusRequestTimeout)
		}

	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}