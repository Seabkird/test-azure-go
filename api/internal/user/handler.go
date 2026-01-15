package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Handler gère les requêtes HTTP pour le domaine User.
type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

// RegisterRoutes permet au handler d'enregistrer ses propres routes sur un routeur parent.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.Search)
	r.Get("/{id}", h.GetByID)
	// r.Put("/{id}", h.Update)
	// r.Delete("/{id}", h.Delete)
}

// =================================================================================
// Handlers HTTP
// =================================================================================

// Create POST /users
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Récupération du tenantID pour faire un return rapide si absent
	tenantID, err := getTenantIDFromContext(ctx)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid tenant context")
		return
	}

	// Décodage du corps JSON vers le DTO d'entrée (CreateUserInput)
	var input CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	defer r.Body.Close()

	// Appel à la couche métier
	newUser, err := h.service.CreateUser(ctx, tenantID, input)
	if err != nil {
		// TODO Ici, vous pourriez vérifier le type d'erreur pour renvoyer 400 ou 409 (conflit)
		// Pour simplifier, on renvoie 500 pour l'instant.
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, newUser)
}

// GetByID gère GET /users/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenantID, err := getTenantIDFromContext(ctx)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Extraction de l'ID depuis l'URL (syntaxe dépendant de votre routeur, ici Chi)
	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing ID parameter")
		return
	}

	// Appel couche métier
	user, err := h.service.GetUser(ctx, tenantID, id)
	if err != nil {
		// TODO: Vérifier si l'erreur est de type "Not Found" pour renvoyer 404
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// Search gère GET /users?nom=...&email=...&limit=10
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenantID, err := getTenantIDFromContext(ctx)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// 1. Parsing des query parameters dans la struct Filter
	filter := parseSearchFilter(r)

	// 2. Appel couche métier
	users, err := h.service.SearchUsers(ctx, tenantID, filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 3. Réponse (Si users est nil, json.Marshal renverra "null", on préfère souvent "[]")
	if users == nil {
		users = []User{}
	}
	respondWithJSON(w, http.StatusOK, users)
}

// =================================================================================
// Helpers privés au Handler (À déplacer potentiellement dans kit/api/http.go)
// =================================================================================

// parseSearchFilter extrait les paramètres d'URL pour construire le filtre.
func parseSearchFilter(r *http.Request) Filter {
	q := r.URL.Query()
	filter := Filter{}

	// Helpers pour parser les string pointers
	if val := q.Get("nom"); val != "" {
		filter.Nom = &val
	}
	if val := q.Get("email"); val != "" {
		filter.Email = &val
	}

	// Pagination avec valeurs par défaut
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20 // Default limit
	}
	filter.Limit = limit

	offset, _ := strconv.Atoi(q.Get("offset"))
	if offset < 0 {
		offset = 0
	}
	filter.Offset = offset

	return filter
}

// respondWithJSON écrit une réponse JSON standard.
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// En prod, gérez l'erreur d'encodage, mais c'est rare qu'elle arrive si payload est valide.
	_ = json.NewEncoder(w).Encode(payload)
}

// respondWithError écrit une réponse d'erreur JSON standard.
func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}

// getTenantIDFromContext est un MOCK.
// Dans votre vrai projet, cela devrait être une fonction exportée de votre package `kit/auth`.
// Elle doit extraire le tenant ID que votre middleware d'authentification a placé dans le contexte.
// TODO: Implémentez cette fonction selon votre logique d'authentification.
func getTenantIDFromContext(ctx context.Context) (string, error) {
	// EXEMPLE FICTIF :
	// tenantID, ok := ctx.Value("tenant_id_key").(string)
	// if !ok || tenantID == "" { return "", fmt.Errorf("no tenant id found") }
	// return tenantID, nil

	// Pour que l'exemple compile, je retourne une valeur fixe.
	// A REMPLACER IMPERATIVEMENT PAR VOTRE LOGIQUE D'AUTH.
	return "tenant_admin", nil
}
