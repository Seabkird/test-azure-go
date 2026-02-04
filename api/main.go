package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"test-api/internal/server"
	"test-api/internal/user"
	"test-api/kit/database/cosmos"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func main() {
	server.InitLogger()

	slog.Info("Démarrage de l'application...")

	// =========================================================================
	// Injection des dépendances
	// =========================================================================

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		slog.Error("Erreur de credential: %v", err)
	}

	endpoint := os.Getenv("COSMOS_ENDPOINT")
	// TODO: Dans un vrai projet, validez que endpoint n'est pas vide
	client, err := azcosmos.NewClient(endpoint, cred, nil)
	if err != nil {
		slog.Error("Erreur création client Cosmos: %v", err)
	}

	userGenericAdapter, err := cosmos.NewAdapter[user.User](client, "TestDB", "UsersContainer")
	if err != nil {
		slog.Error("Impossible d'initialiser l'adaptateur Cosmos pour User: %v", err)
	}

	// TODO on peut si besoin rajouter une petite methode setup dans le domaine user pour garder le main propre
	// ( A voir si on fait un fichier spécifique ou bien juste rajouter dans le handler)
	userRepo := user.NewCosmosRepository(userGenericAdapter)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// =========================================================================
	// Configuration du Routeur HTTP (Chi)
	// =========================================================================

	httpHandler := server.NewRouter(userHandler)

	// =========================================================================
	// Configuration et démarrage du serveur
	// =========================================================================
	port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      httpHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Serveur lancé sur http://localhost:" + port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Erreur serveur : ", err)
	}
}
