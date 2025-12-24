package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"test-api/kit/database"
	"test-api/kit/database/cosmos"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

// App contient toutes les dépendances partagées
type App struct {
	userRepo database.Repository[User]
}

func main() {
	// 1. Initialisation (Config, DB, etc.)
	cred, _ := azcosmos.NewKeyCredential(os.Getenv("COSMOS_KEY"))
	client, _ := azcosmos.NewClientWithKey(os.Getenv("COSMOS_ENDPOINT"), cred, nil)
	repo, _ := cosmos.NewAdapter[User](client, "DB", "Container")

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
		Handler: app.Routes(), // On appelle notre fonction de routing
	}

	fmt.Println("Serveur lancé sur http://localhost:" + port)

	// 4. Lancement
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Erreur serveur : ", err)
	}
}
