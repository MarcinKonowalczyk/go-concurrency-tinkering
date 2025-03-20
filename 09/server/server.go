package main

import (
	"fmt"
	"go-concurrency-tinkering/08/log"

	_log "log"
	"net/http"
	"time"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Println(ctx, "Received request")
		// time.Sleep(2 * time.Second) // Simulate a delay
		// fmt.Fprintln(w, "Hello from the Server!")

		select {
		case <-time.After(2 * time.Second):
			fmt.Fprintln(w, "Hello from the Server!")
		case <-ctx.Done():
			err := ctx.Err()
			log.Println(ctx, "Request cancelled:", err)
			http.Error(w, "Request cancelled", http.StatusRequestTimeout)
		}

	}
	handler = log.Decorate(handler)
	http.HandleFunc("/", handler)

	_log.Println("Starting server on :8080")
	_log.Fatal(http.ListenAndServe(":8080", nil))

}
