package logger

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// Clé pour stocker le logger dans le contexte
type ctxKey struct{}

// Init configure le logger JSON par défaut au démarrage
func Init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		// AddSource: true, // Décommente si tu veux le nom du fichier et la ligne dans les logs
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}

// Middleware : C'est lui qui injecte l'OperationID d'Azure
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Récupération de l'ID Azure
		traceID := ""
		traceParent := r.Header.Get("traceparent")
		if traceParent != "" {
			parts := strings.Split(traceParent, "-")
			if len(parts) > 1 {
				traceID = parts[1]
			}
		}

		// 2. Création d'un logger enrichi
		var l *slog.Logger
		if traceID != "" {
			l = slog.Default().With("operation_Id", traceID)
		} else {
			l = slog.Default()
		}

		// 3. Injection dans le contexte
		ctx := context.WithValue(r.Context(), ctxKey{}, l)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// =============================================================================
// API PROPRE (Helpers)
// =============================================================================

// Info logge un message au niveau INFO en récupérant l'ID du contexte
func Info(ctx context.Context, msg string, args ...any) {
	getLogger(ctx).Info(msg, args...)
}

// Error logge un message au niveau ERROR en récupérant l'ID du contexte
func Error(ctx context.Context, msg string, args ...any) {
	getLogger(ctx).Error(msg, args...)
}

// Warn logge un message au niveau WARN
func Warn(ctx context.Context, msg string, args ...any) {
	getLogger(ctx).Warn(msg, args...)
}

// getLogger récupère le logger du contexte ou renvoie le défaut
func getLogger(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

/*
TODO si on veut utiliser quelque chose de plus propre
dans main
import "test-api/kit/logger" // Importe ton nouveau package

func main() {
    logger.Init() // Juste une fois au début
    // ...
}

dans router.go:
r.Use(logger.Middleware) // Ajoute le middleware

Dans tes Handlers (Health, Users, etc.)
func HealthCheck(w http.ResponseWriter, r *http.Request) {
    // UTILISATION SIMPLE :
    logger.Info(r.Context(), "Health check déclenché")

    // AVEC DES ATTRIBUTS :
    logger.Info(r.Context(), "Vérification DB",
        "status", "OK",
        "latence_ms", 12,
    )

    // EN CAS D'ERREUR :
    if err != nil {
        logger.Error(r.Context(), "Problème critique", "error", err)
    }

    w.Write([]byte("OK"))
}
*/
