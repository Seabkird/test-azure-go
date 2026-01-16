// internal/user/user_test.go

package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"test-api/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =====================================================================================
// TEST DE COMPOSANT (Handler -> Service -> Fake Repo)
// =====================================================================================

func TestCreateUser_Flow(t *testing.T) {
	// 1. SETUP
	fakeRepo := newFakeUserRepository()
	svc := user.NewService(fakeRepo)
	handler := user.NewHandler(svc)

	// Données de test
	const testTenantID = "tenant-123"
	emailToCreate := "arthur@kaamelott.com"

	reqBody := map[string]string{
		"email":  emailToCreate,
		"nom":    "Pendragon",
		"prenom": "Arthur",
	}
	reqBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), user.TenantIDContextKey, testTenantID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// 3. ACTION
	handler.Create(rr, req)

	// 4. ASSERTIONS HTTP
	resp := rr.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var respUser user.User
	err := json.NewDecoder(resp.Body).Decode(&respUser)
	require.NoError(t, err)

	assert.NotEmpty(t, respUser.ID, "L'ID doit être généré")
	assert.Equal(t, emailToCreate, respUser.Email)
	// On vérifie que le tenantID est bien conservé
	assert.Equal(t, testTenantID, respUser.TenantID)

	// 5. ASSERTIONS SUR L'ÉTAT (Fake DB)
	fakeRepo.mu.RLock()
	defer fakeRepo.mu.RUnlock()

	require.Equal(t, 1, len(fakeRepo.data))

	// On reconstruit la clé composée pour vérifier la présence
	storageKey := makeKey(testTenantID, respUser.ID)
	storedUser, exists := fakeRepo.data[storageKey]

	require.True(t, exists, "L'utilisateur doit être trouvé avec la clé tenant#id")
	assert.Equal(t, emailToCreate, storedUser.Email)
}
func TestGetUser_Flow(t *testing.T) {
	// 1. SETUP
	fakeRepo := newFakeUserRepository()
	svc := user.NewService(fakeRepo)
	handler := user.NewHandler(svc)

	const testTenantID = "tenant-456"
	testUserID := uuid.NewString()
	expectedEmail := "leodagan@kaamelott.com"

	// PRÉ-POPULATION
	existingUser := user.User{
		ID:       testUserID,
		TenantID: testTenantID,
		Email:    expectedEmail,
		Nom:      "De Carmelide",
		Prenom:   "Léodagan",
	}
	fakeRepo.data[makeKey(testTenantID, testUserID)] = existingUser

	// --- CORRECTION ICI : ON SETUP UN ROUTEUR CHI ---
	// On crée un mini-routeur juste pour ce test afin que chi.URLParam fonctionne.
	r := chi.NewRouter()
	// On enregistre le handler sur une route avec le paramètre {id}
	r.Get("/users/{id}", handler.GetByID)
	// ------------------------------------------------

	// 2. PRÉPARATION DE LA REQUÊTE
	// L'URL doit correspondre au pattern défini ci-dessus
	targetURL := fmt.Sprintf("/users/%s", testUserID)
	req := httptest.NewRequest(http.MethodGet, targetURL, nil)

	// Injection du TenantID dans le contexte (simulation du middleware auth)
	ctx := context.WithValue(req.Context(), user.TenantIDContextKey, testTenantID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// 3. ACTION
	// --- CORRECTION ICI : ON PASSE PAR LE ROUTEUR ---
	// Au lieu d'appeler handler.GetByID directement, on laisse le routeur faire son travail.
	r.ServeHTTP(rr, req)
	// -----------------------------------------------

	// 4. ASSERTIONS HTTP
	resp := rr.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Le statut devrait être 200 OK")

	var respUser user.User
	err := json.NewDecoder(resp.Body).Decode(&respUser)
	require.NoError(t, err, "Le JSON de réponse doit être valide")

	assert.Equal(t, testUserID, respUser.ID)
	assert.Equal(t, testTenantID, respUser.TenantID)
	assert.Equal(t, expectedEmail, respUser.Email)
}

// =====================================================================================
// IMPLEMENTATION DU FAKE REPOSITORY (COMPATIBLE MULTI-TENANT)
// =====================================================================================

type fakeUserRepository struct {
	mu sync.RWMutex
	// La clé de la map est une combinaison "tenantID#userID"
	data map[string]user.User
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		data: make(map[string]user.User),
	}
}

// Helper pour créer la clé composée
func makeKey(tenantID, id string) string {
	return tenantID + "#" + id
}

// --- Implémentation de l'interface ---

func (f *fakeUserRepository) Create(ctx context.Context, u *user.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if u.TenantID == "" {
		return fmt.Errorf("tenantID is required for creation")
	}

	if u.ID == "" {
		u.ID = uuid.NewString()
	}

	key := makeKey(u.TenantID, u.ID)

	if _, exists := f.data[key]; exists {
		return fmt.Errorf("user already exists with this ID in this tenant")
	}

	f.data[key] = *u

	return nil
}

func (f *fakeUserRepository) GetByID(ctx context.Context, tenantID string, id string) (*user.User, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	key := makeKey(tenantID, id)
	val, ok := f.data[key]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	copyVal := val
	return &copyVal, nil
}

func (f *fakeUserRepository) Update(ctx context.Context, u *user.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if u.TenantID == "" || u.ID == "" {
		return fmt.Errorf("missing tenantID or ID for update")
	}

	key := makeKey(u.TenantID, u.ID)

	if _, exists := f.data[key]; !exists {
		return fmt.Errorf("user not found for update")
	}

	f.data[key] = *u
	return nil
}

func (f *fakeUserRepository) Delete(ctx context.Context, tenantID string, id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	key := makeKey(tenantID, id)
	if _, exists := f.data[key]; !exists {
		return fmt.Errorf("user not found for delete")
	}

	delete(f.data, key)
	return nil
}

// Search implémente une recherche basique en mémoire
func (f *fakeUserRepository) Search(ctx context.Context, tenantID string, filter user.Filter) ([]user.User, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var results []user.User

	// On itère sur toute la map (pas optimal en prod, ok pour un fake)
	for _, v := range f.data {
		// 1. Filtrage obligatoire sur le TenantID
		if v.TenantID != tenantID {
			continue
		}

		// 2. Application des autres filtres de la struct 'Filter'
		// EXEMPLE (à adapter selon votre struct Filter réelle) :
		// if filter.Email != "" && v.Email != filter.Email { continue }

		// Si un filtre est trop complexe pour le fake on peut l'ingorer ou panique
		//panic("Fake repo does not support complex date range filtering. Use integration test instead.")
		results = append(results, v)
	}

	return results, nil
}
