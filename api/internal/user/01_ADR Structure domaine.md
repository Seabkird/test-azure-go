# ADR: [Structure domaine]
**Date:** 15/01/2026

## Context
L'objectif était d'avoir quelque chose de modulaire et facilement duplicable pour pouvoir l'adapter aux différents besoins.
Nous cherchons un équilibre entre une séparation claire des responsabilités (DDD/Clean Arch) et une facilité de navigation dans le code (simplicité).

## Decision
Nous adoptons une structure de Monolithe Modulaire. Chaque module métier (Bounded Context) est autonome dans son dossier.

`internal/user/`
├── user.go    // CONTRAT : Structs principales (DTO/Entities) et Interfaces (Ports)
├── service.go // MÉTIER : Logique d'application (implémente les interfaces de service)
├── handler.go // TRANSPORT HTTP : Gestion des requêtes/réponses et Routage (implémente les interfaces de handler)
└── cosmos.go  // INFRASTRUCTURE : Implémentation DB spécifique (implémente les interfaces de repository)

> Nous avons choisi de regrouper les interfaces dans `user.go` (bien que peu idiomatique en Go pur) pour servir de "sommaire" du domaine.

> Nous avons décidé de **maintenir le routage HTTP (`RegisterRoutes`) dans `handler.go`**. Bien que le déplacer dans `user.go` aurait centralisé la "vue" du module, cela aurait violé la séparation des couches en introduisant des dépendances HTTP dans le fichier de définition du domaine pur.

## Consequences
- `user.go` est la documentation vivante du **métier** et des contrats.
- `handler.go` est la documentation vivante de l'**API HTTP**.
- La séparation est nette : le domaine ignore comment il est exposé sur le web.

## Alternatives Considered

 user/            // Tout ce qui concerne l'utilisateur est ISOLÉ ici
    │   │
    │   ├── model.go     // La structure de données centrale (struct UserEntity...)
    │   │
    │   ├── repository.go // L'INTERFACE définissant le contrat (ex: UserRepository interface { Search(...) })
    │   │
    │   ├── cosmos_repository.go // L'IMPLÉMENTATION CONCRÈTE
    │   │   // C'est LUI qui contient une instance de kit.cosmos.Adapter[UserEntity]
    │   │   // C'est LUI qui contient la méthode SearchSpecificUser avec le SQL WHERE c.age > ...
    │   │
    │   ├── service.go   // COUCHE MÉTIER (Business Logic)
    │   │   // Appelle le repository, valide les données, hache les mots de passe, etc.
    │   │
    │   └── handler.go   // COUCHE HTTP
    │       // Reçoit la requête, décode le JSON entrant, appelle le service, formate la réponse.



