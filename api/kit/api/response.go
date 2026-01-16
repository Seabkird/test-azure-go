package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// errorResponse est la structure JSON standard pour nos erreurs destinées au client.
type errorResponse struct {
	Error string `json:"error"`
}

// respondWithJSON écrit une réponse JSON standard (Statut 2xx).
func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("ERREUR CRITIQUE: échec encodage réponse JSON: %v", err)
	}
}

// respondWithError est le point central de gestion des erreurs.
// Il reçoit l'erreur brute (wrappée), la loggue, et décide de la réponse HTTP.
func RespondWithError(w http.ResponseWriter, err error) {
	// 1. CENTRALISATION DU LOG
	// On loggue ici l'erreur complète avec toute sa trace (%w).
	// C'est ce log que TU regarderas pour debugger.
	// Note: Pour l'instant on utilise le "log" standard. Plus tard, tu injecteras un logger structuré (ex: zap/slog).
	log.Printf("[SERVER ERROR] %v", err)

	// 2. DÉTERMINATION DU STATUS ET DU MESSAGE PUBLIC
	// Par défaut, si on ne reconnaît pas l'erreur, c'est un problème interne (500)
	// et on ne fuite pas les détails techniques au client.
	statusCode := http.StatusInternalServerError
	publicMessage := err.Error()

	// --- ICI viendra plus tard la logique de détection ---
	// Exemple futur (ne pas implémenter maintenant) :
	// if errors.Is(err, kit.ErrNotFound) {
	//     statusCode = http.StatusNotFound
	//     publicMessage = "Resource not found"
	// } else if errors.Is(err, kit.ErrInvalidInput) {
	// ...
	// ----------------------------------------------------

	// 3. ENVOI DE LA RÉPONSE
	RespondWithJSON(w, statusCode, errorResponse{
		Error: publicMessage,
	})
}

// pour plus tard :
// func respondWithError(w http.ResponseWriter, err error) {
// 	// 1. LOG CENTRALISÉ (Toujours l'erreur complète pour le dev)
// 	// Exemple de log généré : "[SERVER ERROR] contexte invalide ou manquant: unauthorized access"
// 	log.Printf("[SERVER ERROR] %v", err)

// 	// 2. DÉTERMINATION DU STATUS ET MESSAGE PUBLIC
// 	statusCode := http.StatusInternalServerError
// 	publicMessage := "Internal Server Error"

// 	// C'est ici que la magie errors.Is opère.
// 	// On vérifie si l'erreur reçue contient une de nos sentinelles.
// 	switch {
// 	case errors.Is(err, kit.ErrUnauthorized):
// 		statusCode = http.StatusUnauthorized
// 		publicMessage = "Accès non autorisé. Veuillez vous authentifier."
// 	case errors.Is(err, kit.ErrNotFound):
// 		statusCode = http.StatusNotFound
// 		publicMessage = "La ressource demandée n'existe pas."
// 	case errors.Is(err, kit.ErrInvalidInput):
// 		statusCode = http.StatusBadRequest
// 		publicMessage = "Les données fournies sont invalides."
// 	// ... Tu peux ajouter d'autres cas ici ...
// 	}

// 	// 3. ENVOI DE LA RÉPONSE
// 	// Dans ton cas, le client recevra un 401 avec {"error": "Accès non autorisé..."}
// 	respondWithJSON(w, statusCode, errorResponse{
// 		Error: publicMessage,
// 	})
// }
