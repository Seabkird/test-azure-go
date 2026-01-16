# ADR: [Structure Globale du Projet - Monolithe Modulaire]

**Date:** 15/01/2026

## Context
Le développement d'un ERP SaaS nécessite une architecture capable de gérer une complexité métier croissante tout en restant maintenable et modulaire. 
Je dois concilier modularité et simplicité opérationnelle sans tomber dans le piège d'un monolithe fortement couplé.

## Decision
J'adopte une structure de **Monolithe Modulaire** respectant les standards Go. 
La logique métier est découpée en modules autonomes (DDD) isolés dans le dossier `internal/`.
Le code technique générique est centralisé dans le dossier `kit/`.

## Consequences
Cette approche garantit une forte cohésion métier et prépare le terrain pour une éventuelle extraction en microservices. 
Elle exige néanmoins une rigueur stricte pour éviter les dépendances cycliques entre les modules métier.

## Alternatives Considered
L'architecture en couches traditionnelle (MVC) a été rejetée car elle ne s'adapte pas bien à la complexité d'un ERP. 
Les microservices ont été écartés pour le moment afin d'éviter une complexité opérationnelle prématurée.

```
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
```