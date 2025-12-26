package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Tes structs de réponse (gardées ici pour l'exemple)
type Reponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	User    User   `json:"user"`
}

type User struct {
	//C'est crucial pour la performance et la sécurité des données entre tes clients.
	TenantID string `json:"tenantId"` // Ta Partition Key (PK)Utilise le TenantID comme clé de partition (Partition Key).
	Id       string `json:"id"`
	Nom      string `json:"nom"`
	Prenom   string `json:"prenom"`
}

func (u User) GetID() string { return u.Nom }

// HandleHome : Ta route texte simple
func HandleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenue sur ton API Go minimaliste test mdr v2!")
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
		Message: "Ceci est une réponse JSON depuis Go",
		User:    u,
	}

	json.NewEncoder(w).Encode(data)
}

func (app *App) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	println("Appel de HandleCreateUser")

	// 1. Décodage du JSON
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Données invalides:"+user.Nom, http.StatusBadRequest)
		return
	}

	// 2. Appel au Repository
	if err := app.userRepo.Create(r.Context(), user); err != nil {
		// En prod, logguez l'erreur réelle ici
		http.Error(w, "Erreur lors de la création"+" "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Réponse
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// HandleGetUser : GET /api/users/{id}?pk=...
func (app *App) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	// 1. Récupération de l'ID via Go 1.22 PathValue
	id := r.PathValue("id")

	// 2. Récupération de la partitionKey (si nécessaire selon ta DB)
	// Exemple: /api/users/123?pk=group1
	pk := r.URL.Query().Get("pk")

	// 3. Appel au Repository
	user, err := app.userRepo.Read(r.Context(), id, pk)
	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
		return
	}

	// 4. Réponse JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// HandleUpdateUser : PUT /api/users/{id}
func (app *App) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Sécurité : on s'assure que l'ID de l'objet correspond à l'URL
	user.Id = id

	if err := app.userRepo.Update(r.Context(), user); err != nil {
		http.Error(w, "Erreur lors de la mise à jour", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// HandleDeleteUser : DELETE /api/users/{id}?pk=...
func (app *App) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	pk := r.URL.Query().Get("pk")
	if pk == "" {
		pk = id
	}

	if err := app.userRepo.Delete(r.Context(), id, pk); err != nil {
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
