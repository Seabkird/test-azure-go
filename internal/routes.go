package main

import "net/http"

func (app *App) Routes() http.Handler {
	// On utilise un ServeMux (le routeur standard de Go)
	mux := http.NewServeMux()

	// Tes routes existantes
	mux.HandleFunc("GET /", HandleHome) // Note: "GET /" marche avec Go 1.22+
	mux.HandleFunc("GET /api/info", HandleInfo)

	// --- Routes CRUD User ---
	mux.HandleFunc("POST /api/user", app.HandleCreateUser)
	mux.HandleFunc("GET /api/user/{id}", app.HandleGetUser)
	mux.HandleFunc("PUT /api/user/{id}", app.HandleUpdateUser)
	mux.HandleFunc("DELETE /api/user/{id}", app.HandleDeleteUser)

	return mux
}
