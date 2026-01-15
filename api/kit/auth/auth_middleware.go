package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// TODO à tester et valider

// Votre clé secrète (doit être la MÊME que celle utilisée pour créer le token lors du login)
// À stocker impérativement dans les Application Settings de la Function App (ex: JWT_SECRET)
var jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

// AuthMiddleware est une fonction qui prend votre handler final et retourne un nouveau handler sécurisé
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Extraire le header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// 2. Le format doit être "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization format. Expected 'Bearer <token>'", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// 3. Parser et valider le token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Vérifier impérativement la méthode de signature (évite les attaques "none")
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Retourner la clé secrète pour la validation
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			// Log l'erreur pour le debug (ne pas renvoyer l'erreur brute au client en prod)
			fmt.Printf("Token validation failed: %v\n", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 4. Optionnel : Extraire des claims (données utilisateur) du token et les passer au contexte
		// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 	 userID := claims["sub"].(string)
		//   // Ajouter userID au contexte de la requête (r.Context()) pour que le handler suivant puisse l'utiliser
		// }

		// 5. Si tout est OK, on appelle la fonction originale
		next.ServeHTTP(w, r)
	}
}
