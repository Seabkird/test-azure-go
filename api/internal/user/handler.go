package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"test-api/kit/api"

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

// RegisterRoutes définit les points d'entrée HTTP pour le module User.
//
// GET /users : Recherche des utilisateurs
// POST /users : Création d'un utilisateur
// GET /users/{id} : Récupération d'un utilisateur par son ID
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
		api.RespondWithError(w, err)
		return
	}

	// Décodage du corps JSON vers le DTO d'entrée (CreateUserInput)
	var input CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		api.RespondWithError(w, err)
		return
	}
	defer r.Body.Close()

	// Appel à la couche métier
	newUser, err := h.service.CreateUser(ctx, tenantID, input)
	if err != nil {
		fmt.Printf("[SERVICE ERROR] Erreur lors de la création de l'utilisateur: %v\n", err)
		// TODO Ici, vous pourriez vérifier le type d'erreur pour renvoyer 400 ou 409 (conflit)
		// Pour simplifier, on renvoie 500 pour l'instant.
		api.RespondWithError(w, err)
		return
	}

	api.RespondWithJSON(w, http.StatusCreated, newUser)
}

// GetByID gère GET /users/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenantID, err := getTenantIDFromContext(ctx)
	if err != nil {
		api.RespondWithError(w, err)
		return
	}

	// Extraction de l'ID depuis l'URL (syntaxe dépendant de votre routeur, ici Chi)
	id := chi.URLParam(r, "id")
	if id == "" {
		api.RespondWithError(w, errors.New("missing id parameter"))
		return
	}

	// Appel couche métier
	user, err := h.service.GetUser(ctx, tenantID, id)
	if err != nil {
		// TODO: Vérifier si l'erreur est de type "Not Found" pour renvoyer 404
		api.RespondWithError(w, err)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, user)
}

// Search gère GET /users?nom=...&email=...&limit=10
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tenantID, err := getTenantIDFromContext(ctx)
	if err != nil {
		api.RespondWithError(w, err)
		return
	}

	// 1. Parsing des query parameters dans la struct Filter
	filter := parseSearchFilter(r)

	// 2. Appel couche métier
	users, err := h.service.SearchUsers(ctx, tenantID, filter)
	if err != nil {
		api.RespondWithError(w, err)
		return
	}

	// 3. Réponse (Si users est nil, json.Marshal renverra "null", on préfère souvent "[]")
	if users == nil {
		users = []User{}
	}
	api.RespondWithJSON(w, http.StatusOK, users)
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

// EXEMPLE FICTIF d'implementation du middleware :
// tenantID, ok := ctx.Value("tenant_id_key").(string)
// if !ok || tenantID == "" { return "", fmt.Errorf("no tenant id found") }
// return tenantID, nil

// Pour que l'exemple compile, je retourne une valeur fixe.
// A REMPLACER IMPERATIVEMENT PAR VOTRE LOGIQUE D'AUTH.

// getTenantIDFromContext est un MOCK.
// Dans votre vrai projet, cela devrait être une fonction exportée de votre package `kit/auth`.
// Elle doit extraire le tenant ID que votre middleware d'authentification a placé dans le contexte.
// TODO: Implémentez cette fonction selon votre logique d'authentification.
func getTenantIDFromContext(ctx context.Context) (string, error) {
	// 1. On essaie de lire la valeur injectée par le test (ou plus tard, le middleware)
	// On utilise la constante publique définie dans user.go
	val := ctx.Value(TenantIDContextKey)

	// 2. Si ça existe et que c'est une string non vide, on l'utilise.
	if tenantID, ok := val.(string); ok && tenantID != "" {
		return tenantID, nil
	}

	// 3. Fallback ou Erreur.
	// Pour le test, si on n'a rien mis dans le contexte, c'est une erreur.
	// Cela forcera le handler à renvoyer une 401 Unauthorized, ce qui est correct.
	return "", fmt.Errorf("no tenant ID found in context (mock implementation)")

	// Note : Si vous voulez garder un comportement par défaut "tenant_admin" pour
	// d'autres tests locaux hors de ce fichier, vous pouvez le remettre ici en fallback,
	// mais cela risque de masquer des erreurs dans vos tests.
	// return "tenant_admin", nil
}
