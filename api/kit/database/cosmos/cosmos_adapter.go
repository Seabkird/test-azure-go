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

func (a *Adapter[T]) Container() *azcosmos.ContainerClient {
	return a.container
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
	pk := azcosmos.NewPartitionKeyString(item.GetTenantID())

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

// TODO à tester et le faire de façon générique car actuellement les filtres sont spécifiques à User
func (a *Adapter[T]) Search(ctx context.Context, filter database.UserFilter, partitionKey string) ([]T, error) {
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	query := "SELECT * FROM c WHERE 1=1" // 1=1 permet d'ajouter des AND facilement
	var params []azcosmos.QueryParameter

	if filter.Category != nil {
		query += " AND c.category = @category"
		params = append(params, azcosmos.QueryParameter{Name: "@category", Value: *filter.Category})
	}

	if filter.MinAge != nil {
		query += " AND c.age >= @minAge"
		params = append(params, azcosmos.QueryParameter{Name: "@minAge", Value: *filter.MinAge})
	}

	if filter.IsActive != nil {
		query += " AND c.isActive = @isActive"
		params = append(params, azcosmos.QueryParameter{Name: "@isActive", Value: *filter.IsActive})
	}

	// Ajout pagination (OFFSET LIMIT est supporté par Cosmos DB moderne)
	if filter.Limit > 0 {
		query += " OFFSET @offset LIMIT @limit"
		params = append(params, azcosmos.QueryParameter{Name: "@offset", Value: filter.Offset})
		params = append(params, azcosmos.QueryParameter{Name: "@limit", Value: filter.Limit})
	}

	// Exécution avec le SDK Cosmos...
	// --- CORRECTION 1 : ORDER BY OBLIGATOIRE ---
	// Cosmos DB exige un ORDER BY pour utiliser OFFSET/LIMIT.
	// Ici je trie par ID, mais ça pourrait être c.createdAt DESC
	if filter.Limit > 0 {
		query += " ORDER BY c.id ASC OFFSET @offset LIMIT @limit"
		params = append(params, azcosmos.QueryParameter{Name: "@offset", Value: filter.Offset})
		params = append(params, azcosmos.QueryParameter{Name: "@limit", Value: filter.Limit})
	}

	// Préparation des options
	queryOptions := azcosmos.QueryOptions{
		QueryParameters: params,
	}

	// --- CORRECTION 2 : Optimisation Partition Key ---
	// Si "Category" est ta partition key, il faut l'ajouter ici pour éviter un scan complet
	// if filter.Category != nil {
	// 	 queryOptions.PartitionKey = azcosmos.NewPartitionKeyString(*filter.Category)
	// }

	// Création du Pager
	pager := a.container.NewQueryItemsPager(query, pk, &queryOptions)

	var results []T

	// --- CORRECTION 3 : Boucle et Désérialisation ---
	for pager.More() {
		// Récupération de la page (appel réseau)
		response, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la requête cosmos: %w", err)
		}

		// Chaque réponse contient une liste d'items sous forme de []byte (JSON brut)
		for _, bytes := range response.Items {
			var item T
			// On transforme le JSON en struct Go T
			if err := json.Unmarshal(bytes, &item); err != nil {
				return nil, fmt.Errorf("erreur de désérialisation: %w", err)
			}
			results = append(results, item)
		}
	}

	return results, nil
}
