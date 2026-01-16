package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"

	"test-api/kit/database/cosmos"
)

// cosmosRepository est l'implémentation spécifique du Repository pour le domaine User.
type cosmosRepository struct {
	genericAdapter  *cosmos.Adapter[User]
	containerClient *azcosmos.ContainerClient
}

func NewCosmosRepository(adapter *cosmos.Adapter[User]) Repository {
	// On suppose que le kit a la méthode Container()
	return &cosmosRepository{
		genericAdapter: adapter,
		// On récupère le client bas niveau depuis l'adaptateur
		containerClient: adapter.Container(),
	}
}

// =================================================================================
// Méthodes CRUD simples (Délégation à l'adapteur générique)
// =================================================================================

func (r *cosmosRepository) Create(ctx context.Context, user *User) error {
	// L'adapteur générique sait comment appeler GetTenantID() grâce à l'interface Entity
	return r.genericAdapter.Create(ctx, *user)
}

func (r *cosmosRepository) GetByID(ctx context.Context, tenantID string, id string) (*User, error) {
	user, err := r.genericAdapter.Read(ctx, id, tenantID)

	if err != nil {
		// Gestion spécifique de l'erreur "not found" de Cosmos pour que ne soit aps considéré comme une erreur technique.
		var responseErr *azcore.ResponseError
		if errors.As(err, &responseErr) {
			if responseErr.StatusCode == 404 {
				return nil, nil
			}
		}
		// Sinon, c'est une vraie erreur technique (timeout, auth, etc.)
		return nil, err
	}

	if user.ID == "" {
		return nil, nil
	}
	return &user, nil
}

func (r *cosmosRepository) Update(ctx context.Context, user *User) error {
	return r.genericAdapter.Update(ctx, *user)
}

func (r *cosmosRepository) Delete(ctx context.Context, tenantID string, id string) error {
	return r.genericAdapter.Delete(ctx, id, tenantID)
}

// =================================================================================
// Méthodes de Recherche Spécifiques (Implémentation directe ici)
// =================================================================================

// Search implémente la recherche multicritères spécifique aux users.
func (r *cosmosRepository) Search(ctx context.Context, tenantID string, filter Filter) ([]User, error) {
	pk := azcosmos.NewPartitionKeyString(tenantID)

	// Construction de la requête SQL de base
	// IMPORTANT : On filtre TOUJOURS par tenantID dans la clause WHERE pour la sécurité.
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("SELECT * FROM c WHERE c.tenantID = @tenantId")

	// Initialisation des paramètres avec le tenantId
	params := []azcosmos.QueryParameter{
		{Name: "@tenantId", Value: tenantID},
	}

	// Ajout dynamique des filtres optionnels
	if filter.Nom != nil {
		queryBuilder.WriteString(" AND c.nom = @nom")
		params = append(params, azcosmos.QueryParameter{Name: "@nom", Value: *filter.Nom})
	}

	if filter.Email != nil {
		queryBuilder.WriteString(" AND c.email = @email")
		params = append(params, azcosmos.QueryParameter{Name: "@email", Value: *filter.Email})
	}

	// Ajout de la pagination (ORDER BY obligatoire pour OFFSET/LIMIT)
	queryBuilder.WriteString(" ORDER BY c._ts DESC")

	if filter.Limit > 0 {
		queryBuilder.WriteString(" OFFSET @offset LIMIT @limit")
		params = append(params, azcosmos.QueryParameter{Name: "@offset", Value: filter.Offset})
		params = append(params, azcosmos.QueryParameter{Name: "@limit", Value: filter.Limit})
	}

	// Préparation des options de requête
	queryOptions := azcosmos.QueryOptions{
		QueryParameters: params,
	}

	// Exécution avec le Pager (via le container client brut)
	// Note : Votre Adapteur générique pourrait exposer une méthode "ExecuteQuery" pour éviter
	// d'avoir accès au containerClient ici, mais pour l'instant, faisons simple.
	pager := r.containerClient.NewQueryItemsPager(queryBuilder.String(), pk, &queryOptions)

	// Boucle de récupération des résultats
	// Vous pouvez réutiliser la méthode privée de votre adapteur si elle existe,
	// sinon on duplique cette boucle standard ici.
	var results []User

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("cosmos query failed: %w", err)
		}

		for _, bytes := range response.Items {
			// Désérialisation via la méthode générique de l'adapteur si possible,
			// sinon on le fait à la main ici. Votre adapteur semble faire le unmarshal
			// à l'intérieur de ses méthodes Read/Search, mais pas via une fonction helper exposée.
			// On le refait donc ici :
			var item User
			if err := json.Unmarshal(bytes, &item); err != nil {
				// Log error, continue or return...
				return nil, fmt.Errorf("failed to unmarshal user json: %w", err)
			}
			results = append(results, item)
		}
	}

	return results, nil
}
