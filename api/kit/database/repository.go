package database

import (
	"context"
)

// Entity définit ce qu'est un objet stockable de base.
// Cosmos DB a besoin d'un champ "id" en minuscule JSON.
type Entity interface {
	GetID() string
}

// Repository définit les opérations CRUD standard.
type Repository[T Entity] interface {
	Create(ctx context.Context, item T) error
	Read(ctx context.Context, id string, partitionKey string) (T, error)
	Update(ctx context.Context, item T) error
	Delete(ctx context.Context, id string, partitionKey string) error
	// Tu pourras ajouter Query() plus tard
}
