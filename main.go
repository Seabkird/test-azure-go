package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Structure pour la réponse JSON
type Reponse struct {
	Status      string      `json:"status"`
	Message     string      `json:"message"`
	Utilisateur Utilisateur `json:"utilisateur"`
}

type Utilisateur struct {
	Nom    string `json:"nom"`
	Prenom string `json:"prenom"`
}

func main() {
	// 1. Route simple (Texte brut)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Bienvenue sur ton API Go minimaliste !")
	})

	// 2. Route API (Retourne du JSON)
	http.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		// On définit le header pour dire "c'est du JSON"
		w.Header().Set("Content-Type", "application/json")

		// On prépare les données
		data := Reponse{
			Status:  "success",
			Message: "Ceci est une réponse JSON depuis Go",
			Utilisateur: Utilisateur{
				Nom:    "de Casteljau",
				Prenom: "Raphaël",
			},
		}

		// On encode et on envoie
		json.NewEncoder(w).Encode(data)
	})

	// 3. Lancement du serveur
	port := ":8080"
	fmt.Println("Serveur lancé sur http://localhost" + port)

	// Bloque le programme et écoute les requêtes
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Erreur lors du lancement du serveur : ", err)
	}
}
