package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Tes structs de réponse (gardées ici pour l'exemple)
type Reponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	User    User   `json:"user"`
}

type User struct {
	Nom    string `json:"nom"`
	Prenom string `json:"prenom"`
}

func (u User) GetID() string { return u.Nom }

// HandleHome : Ta route texte simple
func (app *App) HandleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenue sur ton API Go minimaliste test mdr v2!")
}

// HandleInfo : Ta route JSON
func (app *App) HandleInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := Reponse{
		Status:  "success",
		Message: "Ceci est une réponse JSON depuis Go",
		User: User{
			Nom:    "de Casteljau",
			Prenom: "Raphaël",
		},
	}

	json.NewEncoder(w).Encode(data)
}
