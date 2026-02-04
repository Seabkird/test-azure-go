package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// InitLogger configure le logger par défaut pour écrire du JSON sur Stdout
func InitLogger() {
	// On utilise JSON, car Azure parse automatiquement le JSON dans Log Analytics
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

// AzureTraceMiddleware est le middleware qui récupère l'ID d'Azure
func AzureTraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := ""

		// Azure Functions envoie l'ID dans le header "traceparent"
		// Format W3C: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
		traceParent := r.Header.Get("traceparent")

		if traceParent != "" {
			parts := strings.Split(traceParent, "-")
			if len(parts) > 1 {
				traceID = parts[1] // C'est l'Operation ID que Application Insights utilise
			}
		}

		// Si on a trouvé un ID, on l'ajoute au logger pour cette requête
		if traceID != "" {
			// On crée un logger avec l'attribut "operation_Id" déjà rempli
			// Note: "operation_Id" est le nom de colonne exact dans App Insights
			ctx := context.WithValue(r.Context(), "logger", slog.With("operation_Id", traceID))
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// Log retourne une instance de logger enrichie avec le contexte (Trace ID)
func Log(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value("logger").(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
