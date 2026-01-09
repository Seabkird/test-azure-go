package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Reponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	User    User   `json:"user"`
}

// HandleHome : Ta route texte simple
func HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"message": "Bienvenue sur ton API Go minimaliste test mdr v2!",
	}

	json.NewEncoder(w).Encode(response)
}

// HandleInfo : Ta route JSON
func HandleInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}

	data := Reponse{
		Status:  "success",
		Message: "Il s'agit de la route /api/health",
		User:    u,
	}

	json.NewEncoder(w).Encode(data)
}
