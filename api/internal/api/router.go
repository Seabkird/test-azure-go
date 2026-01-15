package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// Importez les packages de vos domaines pour accéder aux handlers
	"test-api/internal/user"
	// "test-api/internal/product" // Futur import
)

// NewRouter initialise le routeur Chi principal.
func NewRouter(userHandler *user.Handler /*, productHandler *product.Handler */) http.Handler {
	r := chi.NewRouter()

	// =========================================================================
	// Middlewares Globaux
	// =========================================================================
	// RequestID ajoute un ID unique à chaque requête (utile pour les logs)
	r.Use(middleware.RequestID)
	// Logger logge les détails de la requête HTTP entrante
	r.Use(middleware.Logger)
	// Recoverer empêche le serveur de planter en cas de panic dans un handler
	r.Use(middleware.Recoverer)

	// Vous pouvez ajouter vos propres middlewares ici (ex: auth globale, CORS)
	// r.Use(myAuthMiddleware)

	// =========================================================================
	// Routes de base
	// =========================================================================
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// =========================================================================
	// Montage des routes API des différents domaines (/api/...)
	// =========================================================================
	// On groupe toutes les routes API sous le préfixe "/api"
	r.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Route("/users", func(userRouter chi.Router) {
			userHandler.RegisterRoutes(userRouter)
		})

		// --- (Futur) Domaine PRODUCT ---
		/*
			apiRouter.Route("/products", func(productRouter chi.Router) {
				productHandler.RegisterRoutes(productRouter)
			})
		*/
	})

	return r
}
