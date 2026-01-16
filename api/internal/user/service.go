package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type serviceImpl struct {
	repo Repository
}

// -- Définition des erreurs métier --

var ErrUserNotFound = fmt.Errorf("user not found")
var ErrEmailAlreadyExists = fmt.Errorf("email already registered for this tenant")

// ErrInvalidInput est une erreur générique de validation.
type ErrInvalidInput struct {
	Field   string
	Message string
}

func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input for field '%s': %s", e.Field, e.Message)
}

// =================================================================================
// Implémentation du Service
// =================================================================================

func NewService(r Repository) Service {
	return &serviceImpl{
		repo: r,
	}
}

func (s *serviceImpl) CreateUser(ctx context.Context, tenantID string, input CreateUserInput) (*User, error) {
	// 1. Nettoyage et validation de base des entrées
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if email == "" {
		return nil, ErrInvalidInput{Field: "email", Message: "cannot be empty"}
	}
	if !strings.Contains(email, "@") {
		// C'est une validation simpliste, utilisez une regex en prod
		return nil, ErrInvalidInput{Field: "email", Message: "invalid format"}
	}

	if strings.TrimSpace(input.Nom) == "" {
		return nil, ErrInvalidInput{Field: "nom", Message: "cannot be empty"}
	}

	// 2. Validation métier : Vérifier l'unicité de l'email dans ce tenant.
	// On doit faire un appel au repo pour vérifier.
	// ATTENTION : C'est une vérification "soft". Entre cet appel et l'insertion,
	// une autre requête concurrente pourrait passer. C'est un compromis classique en NoSQL.
	// La vraie unicité doit être gérée par la DB si possible (index unique composite tenantID+email sur Cosmos).
	checkFilter := Filter{Email: &email, Limit: 1}
	existingUsers, err := s.repo.Search(ctx, tenantID, checkFilter)
	if err != nil {
		// Si erreur technique DB, on remonte.
		return nil, fmt.Errorf("failed to check existing email: %w", err)
	}
	if len(existingUsers) > 0 {
		// Règle métier violée
		return nil, ErrEmailAlreadyExists
	}

	// 3. Enrichissement des données et création de l'entité finale
	newUser := &User{
		TenantID: tenantID,            // Vient du contexte de sécurité
		ID:       uuid.New().String(), // Génération de l'ID technique
		Email:    email,
		Nom:      strings.TrimSpace(input.Nom),
		Prenom:   strings.TrimSpace(input.Prenom),
		// Ici on ajouterait:
		// CreatedAt: time.Now().UTC(),
		// IsActive:  true,
	}

	// 4. Persistance via le repository
	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user in repo: %w", err)
	}

	// 5. Retourner l'entité créée
	return newUser, nil
}

// GetUser implémente la logique de récupération simple.
func (s *serviceImpl) GetUser(ctx context.Context, tenantID string, id string) (*User, error) {
	// Validation simple de l'ID
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidInput{Field: "id", Message: "invalid UUID format"}
	}

	user, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		// Ici, il faudrait que le repo retourne une erreur standardisée si not found.
		// Si le repo retourne nil, nil, on peut gérer comme ça :
		// if user == nil { return nil, ErrUserNotFound }
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	// IMPORTANT : Sécurité multi-tenant
	// Même si le repo est censé filtrer par TenantID, c'est une bonne pratique de défense en profondeur
	// de revérifier ici que l'objet retourné appartient bien au tenant courant.
	if user.TenantID != tenantID {
		// Cela ne devrait théoriquement pas arriver si le repo fait son travail,
		// mais si ça arrive, c'est une faille de sécurité critique.
		// On loggue une erreur critique et on dit "Not Found".
		// log.Critical("Security alert: cross-tenant access attempted detected in service layer")
		return nil, ErrUserNotFound
	}

	return user, nil
}

// SearchUsers implémente la logique de recherche.
func (s *serviceImpl) SearchUsers(ctx context.Context, tenantID string, filter Filter) ([]User, error) {
	// Ici, on pourrait appliquer des règles métier sur le filtre.
	// Par exemple, forcer une limite max si elle n'est pas fournie pour éviter de tuer la DB.
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100 // Hard cap métier
	}

	// Appel au repository
	users, err := s.repo.Search(ctx, tenantID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	// Pas besoin de défense en profondeur sur le tenantID ici car le repo est censé
	// appliquer la clause "WHERE tenantId = X" sur toute la liste.

	return users, nil
}

func (s *serviceImpl) UpdateUser(ctx context.Context, tenantID string, id string, input UpdateUserInput) (*User, error) {
	// TODO: Implémenter
	// 1. GetUser(id) pour vérifier qu'il existe et qu'on est dans le bon tenant.
	// 2. Appliquer les modifications des champs non-nil de l'input sur l'entité récupérée.
	// 3. repo.Update(userUpdated)
	return nil, fmt.Errorf("not implemented")
}

func (s *serviceImpl) DeleteUser(ctx context.Context, tenantID string, id string) error {
	// TODO: Implémenter
	// Validation UUID
	// repo.Delete(...)
	return fmt.Errorf("not implemented")
}
