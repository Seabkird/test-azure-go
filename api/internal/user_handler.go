package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// HandleCreateUser : POST /api/users
func (app *App) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	println("Appel de HandleCreateUser")

	// 1. Décodage du JSON
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Données invalides:"+user.Nom, http.StatusBadRequest)
		return
	}

	user.Id = uuid.New().String()
	// TODO: Récupérer le TenantID depuis le contexte/authentification
	user.TenantID = "tenant_admin"

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
	// 1. Récupération de l'id via le path
	id := r.PathValue("id")

	// 2. Récupération de la partitionKey via le query param
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
		log.Fatalf("Erreur création client Cosmos: %v", err)
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
