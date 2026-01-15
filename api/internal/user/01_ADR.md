Ici on fait du DDD  donc on mets tout ce qui est lié au domaine utilisateur.
On l'appelle souvent Monolithe Modulaire (Modular Monolith) structuré selon les principes du Domain-Driven Design (DDD) (ou du moins, inspiré par).

internal/user/
├── user.go // Contient souvent la struct principale 'User' et les interfaces clés
├── service.go
├── handler.go
└── cosmos.go // Pour l'implémentation DB spécifique

> Je suis parti sur cette solution car c'est la plus simple et condensé possible tout en conservé bien la sépération des différents éléments. 



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



