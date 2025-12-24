package main

import "net/http"

func (app *App) Routes() http.Handler {
	// On utilise un ServeMux (le routeur standard de Go)
	mux := http.NewServeMux()

	// Tes routes existantes
	mux.HandleFunc("GET /", app.HandleHome) // Note: "GET /" marche avec Go 1.22+
	mux.HandleFunc("GET /api/info", app.HandleInfo)

	// Exemple futur pour ton CRUD User (quand tu seras prêt)
	//mux.HandleFunc("POST /api/users", app.HandleCreateUser)

	return mux
}
