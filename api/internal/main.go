package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"test-api/kit/database"
	"test-api/kit/database/cosmos"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// App contient toutes les dépendances partagées
type App struct {
	userRepo database.Repository[User]
}

func main() {
	// 1. Initialisation (Config, DB, etc.)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Erreur de credential: %v", err)
	}

	endpoint := os.Getenv("COSMOS_ENDPOINT")
	client, err := azcosmos.NewClient(endpoint, cred, nil)
	if err != nil {
		log.Fatalf("Erreur création client Cosmos: %v", err)
	}

	repo, _ := cosmos.NewAdapter[User](client, "TestDB", "NomContainer")

	// 2. Création de l'application avec ses dépendances
	app := &App{
		userRepo: repo,
	}

	// 3. Configuration du serveur
	port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: app.Routes(),
	}

	fmt.Println("Serveur lancé sur http://localhost:" + port)

	// 4. Lancement
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Erreur serveur : ", err)
	}
}
