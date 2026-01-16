Ne génère pas de code tant que je ne te l'ai pas demandé et soit relativement court dans tes réponse. 
N'hésite pas à me contredire si jamais je pars sur une mauvaise solution/piste. 
N'hésite pas à me demander si tu as besoin de plus d'information. 

Voici ma structure pour un projet de saas erp en tant que dev solo :
.
├── go.mod
├── go.sum
├── main.go              // Point d'entrée principal (setup serveur HTTP, injection dépendances)
│
├── kit/                 // "Toolkit" : Code technique réutilisable, agnostique du métier
│   ├── api/             // Outils API génériques
│   │   └── response.go  // Helpers standardisés (ex: RespondWithJSON, RespondWithError)
│   ├── auth/
│   │   └── middleware.go
│   └── database/
│       ├── cosmos/
│       │   // ADAPTEUR GÉNÉRIQUE 
│       │   // Il contient les méthodes CRUD de base (GetByID, Create, Update...)
│       │   // Il ne connaît PAS les champs spécifiques (age, category...)
│       │   └── cosmos_adapter.go
│       └── errors.go    // Erreurs standard de base de données (ex: ErrNotFound)
│
└── internal/            // Cœur de votre application métier (non importable de l'extérieur)
    │
    ├── server/          // Couche HTTP globale de haut niveau
    │   └── router.go    // Configuration des routes Chi (l'orchestrateur du serveur)
    │
    │   // --- DOMAINE : USER ---
    ├── user/
        ├── user.go // Contient la struct principale 'User' et les interfaces clés
        ├── service.go  // Code métier
        ├── handler.go  // Port http
        └── cosmos.go // Pour l'implémentation DB spécifique
    │   // --- DOMAINE : FUTURE FEATURE (ex: PRODUCT) ---
    ├── product/
    │   ├── model.go
    │   ├── cosmos_repository.go // Contiendra sa propre méthode SearchProduct spécifique
    │   └── ...
    │
    └── config/          // Chargement de la configuration (env vars)
        └── config.go