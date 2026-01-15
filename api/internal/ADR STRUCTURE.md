.
├── go.mod
├── go.sum
├── main.go              // Point d'entrée principal (setup serveur HTTP, injection dépendances)
│
├── kit/                 // "Toolkit" : Code technique réutilisable, agnostique du métier
│   ├── auth/
│   │   └── middleware.go
│   └── database/
│       ├── cosmos/
│       │   // VOTRE ADAPTEUR GÉNÉRIQUE EST ICI
│       │   // Il contient les méthodes CRUD de base (GetByID, Create, Update...)
│       │   // Il ne connaît PAS les champs spécifiques (age, category...)
│       │   └── cosmos_adapter.go 
│       └── errors.go    // Erreurs standard de base de données (ex: ErrNotFound)
│
└── internal/            // Cœur de votre application métier (non importable de l'extérieur)
    │
    ├── api/             // Couche HTTP globale
    │   ├── router.go    // Configuration des routes (Chi, Gin, Echo...)
    │   └── response.go  // Helpers pour standardiser les réponses JSON (ex: respondWithJSON, respondWithError)
    │
    │   // --- DOMAINE : USER ---
    ├── internal/user/
        ├── user.go // Contient souvent la struct principale 'User' et les interfaces clés
        ├── service.go
        ├── handler.go
        └── cosmos.go // Pour l'implémentation DB spécifique
    │   // --- DOMAINE : FUTURE FEATURE (ex: PRODUCT) ---
    ├── product/
    │   ├── model.go
    │   ├── cosmos_repository.go // Contiendra sa propre méthode SearchProduct spécifique
    │   └── ...
    │
    └── config/          // Chargement de la configuration (env vars)
        └── config.go