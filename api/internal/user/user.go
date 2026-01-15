package user

import (
	"context"
)

// =================================================================================
// Modèles de Données (Entities & DTOs)
// =================================================================================

type User struct {
	// TenantID est la clé de partition (critique pour Cosmos DB dans un SaaS multi-tenant).
	TenantID string `json:"tenantID"`

	// ID est l'identifiant unique de l'utilisateur (UUID).
	ID string `json:"id"`

	Email  string `json:"email"`
	Nom    string `json:"nom"`
	Prenom string `json:"prenom"`

	// Vous ajouterez sûrement ici plus tard :
	// HashedPassword string `json:"-"` // Le "-" évite de le renvoyer dans le JSON
	// CreatedAt      time.Time `json:"createdAt"`
	// IsActive       bool      `json:"isActive"`
}

// GetID retourne l'identifiant unique.
func (u User) GetID() string {
	return u.ID
}

// GetTenantID retourne la clé de partition.
func (u User) GetTenantID() string {
	return u.TenantID
}

// ---------------------------------------------------------------------------------
// Modèles d'Entrée (Input DTOs)
// Ces structures servent à valider les données entrant dans votre API.
// On ne veut pas exposer la struct User complète lors de la création (on ne veut pas que l'utilisateur choisisse son ID ou son TenantID).
// ---------------------------------------------------------------------------------

// CreateUserInput définit les données nécessaires pour créer un nouvel utilisateur.
type CreateUserInput struct {
	Email  string `json:"email"`
	Nom    string `json:"nom"`
	Prenom string `json:"prenom"`
	// Password string `json:"password"`
}

// UpdateUserInput définit les champs modifiables d'un utilisateur.
// L'utilisation de pointeurs (*) permet de savoir si un champ a été fourni ou non (pour faire du PATCH).
type UpdateUserInput struct {
	Email  *string `json:"email,omitempty"`
	Nom    *string `json:"nom,omitempty"`
	Prenom *string `json:"prenom,omitempty"`
}

// Filter définit les critères de recherche pour la méthode Search.
type Filter struct {
	// Pointeurs pour distinguer la recherche d'une chaîne vide vs pas de filtre sur ce champ
	Email *string
	Nom   *string

	// Pagination
	Offset int
	Limit  int
}

// =================================================================================
// Interfaces (Contrats)
// =================================================================================

// Service définit le contrat de la couche métier (Business Logic).
type Service interface {
	CreateUser(ctx context.Context, tenantID string, input CreateUserInput) (*User, error)
	GetUser(ctx context.Context, tenantID string, id string) (*User, error)
	UpdateUser(ctx context.Context, tenantID string, id string, input UpdateUserInput) (*User, error)
	DeleteUser(ctx context.Context, tenantID string, id string) error

	SearchUsers(ctx context.Context, tenantID string, filter Filter) ([]User, error)
}

// Repository définit le contrat pour la couche de persistance (Base de données).
type Repository interface {
	// Create persiste un utilisateur complet.
	Create(ctx context.Context, user *User) error

	// GetByID récupère un utilisateur par sa clé composite (TenantID + ID).
	GetByID(ctx context.Context, tenantID string, id string) (*User, error)

	// Search exécute une recherche multi-critères.
	Search(ctx context.Context, tenantID string, filter Filter) ([]User, error)

	// Update met à jour un utilisateur existant.
	Update(ctx context.Context, user *User) error

	// Delete supprime un utilisateur.
	Delete(ctx context.Context, tenantID string, id string) error
}
