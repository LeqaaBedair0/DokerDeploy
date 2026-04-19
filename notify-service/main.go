package main

import (
	"io"
	"log"
	"net/http"
)

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Simulate sending an email or notification
	log.Printf("📩 EVENT RECEIVED | New Task Notification: %s\n", string(body))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification logged successfully"))
}

func main() {
	http.HandleFunc("/notify", notifyHandler)
	log.Println("🚀 Notify Service running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
