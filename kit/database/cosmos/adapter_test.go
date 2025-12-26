package cosmos_test

import (
	"context"
	"os"
	"testing"
	"time"

	"test-api/kit/database/cosmos"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// Struct de test locale
type TestUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (u TestUser) GetID() string { return u.ID }

func TestCosmosAdapter_Integration(t *testing.T) {
	// 1. Vérifier si on doit lancer le test
	endpoint := os.Getenv("COSMOS_ENDPOINT")
	key := os.Getenv("COSMOS_KEY")
	if endpoint == "" || key == "" {
		t.Skip("Skipping integration test: COSMOS_ENDPOINT or COSMOS_KEY not set")
	}

	// 2. Setup
	cred, _ := azcosmos.NewKeyCredential(key)
	client, _ := azcosmos.NewClientWithKey(endpoint, cred, nil)

	// Utilise une base/container de test dédié si possible
	repo, err := cosmos.NewAdapter[TestUser](client, "TestDB", "TestContainer")
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := TestUser{ID: "integration-test-1", Email: "test@example.com"}

	// 3. Scénario de test : Create -> Get -> Delete

	// CREATE
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// GET
	got, err := repo.Read(ctx, user.ID, user.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, got.Email)
	}

	// DELETE
	if err := repo.Delete(ctx, user.ID, user.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// VERIFY DELETE (Should fail)
	_, err = repo.Read(ctx, user.ID, user.ID)
	if err == nil {
		t.Error("Expected error after delete, got nil")
	}
}
