package main

import (
	"encoding/json"
	"github.com/edpuzn/HackAthon/cmd"
	"log"
	"net/http"
)

type RequestBody struct {
	Query string `json:"query"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// CORS başlıklarını ekle
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Preflight (OPTIONS) isteği yönetimi
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Kullanıcı kimliğini belirle
	userID := "exampleUserID" // Örnek bir kullanıcı kimliği kullanıyoruz

	// GetResponse yerine HandleAPIRequest'i kullanıyoruz
	response := cmd.HandleAPIRequest(userID, reqBody.Query)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/search", searchHandler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
