package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"test-api/internal/user"
	"test-api/kit/logger"
)

func NewRouter(userHandler *user.Handler) http.Handler {
	r := chi.NewRouter()

	// =========================================================================
	// Middlewares Globaux
	// =========================================================================
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// 2. On utilise le Middleware de notre nouveau package "logger"
	// (Remplace "middleware.Logger" de Chi qui fait des logs texte moches)
	r.Use(logger.Middleware)

	// =========================================================================
	// Routes de base
	// =========================================================================
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// 3. UTILISATION PROPRE : On appelle le package logger importé
		logger.Info(r.Context(), "Health check déclenché : tout va bien")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// =========================================================================
	// Montage des routes API des différents domaines (/api/...)
	// =========================================================================
	// On groupe toutes les routes API sous le préfixe "/api"
	r.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Use(devTenantMiddleware)

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

// --- AJOUT TEMPORAIRE : Middleware pour simuler un tenant ---
// Ce middleware injecte un ID "en dur" pour le développement.
// À RETIRER une fois le vrai middleware d'auth implémenté dans kit/auth/.
func devTenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. La valeur en dur que tu veux utiliser pour tes tests actuels
		hardcodedTenantID := "tenant-admin"

		// 2. On crée un nouveau contexte dérivé de celui de la requête,
		// en y ajoutant la valeur avec la clé définie dans kit/contextkeys
		ctx := context.WithValue(r.Context(), user.TenantIDContextKey, hardcodedTenantID)

		// 3. Créer une nouvelle requête avec ce nouveau contexte
		rWithCtx := r.WithContext(ctx)

		// 4. Passer la main au handler suivant avec la requête modifiée
		next.ServeHTTP(w, rWithCtx)
	})
}
