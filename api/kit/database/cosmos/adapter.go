package cosmos

import (
	"context"
	"encoding/json"
	"fmt"

	"test-api/kit/database"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// Adapter implémente database.Repository pour Cosmos DB.
type Adapter[T database.Entity] struct {
	container *azcosmos.ContainerClient
}

// NewAdapter crée une nouvelle instance du repository.
func NewAdapter[T database.Entity](client *azcosmos.Client, dbName, containerName string) (*Adapter[T], error) {
	db, err := client.NewDatabase(dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to get database client: %w", err)
	}

	container, err := db.NewContainer(containerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get container client: %w", err)
	}

	return &Adapter[T]{
		container: container,
	}, nil
}

func (a *Adapter[T]) Create(ctx context.Context, item T) error {
	pk := azcosmos.NewPartitionKeyString(item.GetTenantID())

	b, err := json.Marshal(item)
	if err != nil {
		return err
	}

	_, err = a.container.CreateItem(ctx, pk, b, nil)
	return err
}

func (a *Adapter[T]) Read(ctx context.Context, id string, partitionKey string) (T, error) {
	var item T
	// Création d'une instance vide pour éviter le nil pointer si T est un pointeur
	// Note: avec les génériques, c'est parfois tricky, l'appelant recevra la zero-value en cas d'erreur.

	pk := azcosmos.NewPartitionKeyString(partitionKey)

	res, err := a.container.ReadItem(ctx, pk, id, nil)
	if err != nil {
		return item, err
	}

	err = json.Unmarshal(res.Value, &item)
	return item, err
}

func (a *Adapter[T]) Update(ctx context.Context, item T) error {
	pk := azcosmos.NewPartitionKeyString(item.GetID())

	b, err := json.Marshal(item)
	if err != nil {
		return err
	}

	// ReplaceItem écrase l'élément existant
	_, err = a.container.ReplaceItem(ctx, pk, item.GetID(), b, nil)
	return err
}

func (a *Adapter[T]) Delete(ctx context.Context, id string, partitionKey string) error {
	pk := azcosmos.NewPartitionKeyString(partitionKey)
	_, err := a.container.DeleteItem(ctx, pk, id, nil)
	return err
}
