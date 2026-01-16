# TODO REST API User

- [ ] **Refactorisation HTTP Helpers** : Déplacer  `getTenantIDFromContext` dans un package partagé (ex: `internal/kit/http`).
- [ ] **Erreurs Structurées** : Créer une struct `APIError` standard (code, message) au lieu d'une simple string.
- [ ] **Sécurité Erreurs 500** : Ne jamais renvoyer l'erreur brute `err.Error()` en prod. Logger l'erreur et renvoyer un message générique au client.
- [ ] **Mapping Erreurs Métier** : Dans les handlers, vérifier le type d'erreur retourné par le service (ex: `ErrUserAlreadyExists`, `ErrNotFound`) pour renvoyer le bon status code (409, 404) au lieu de 500 systématiquement.
- [ ] **Améliorer Search Parsing** : Automatiser le parsing des paramètres d'URL dans `parseSearchFilter` (envisager la lib `github.com/go-playground/form`) pour éviter le code répétitif.

Rajouter ces test la qui semble indispensable : 

Tes tests actuels sont d'excellents tests de composants "happy path". Ils valident que la chaîne entière (Handler HTTP -> Service -> Fake DB) fonctionne quand tout va bien. Pour un solopreneur, c'est le type de test le plus rentable.
Cependant, il manque des tests cruciaux pour garantir la robustesse de ton API, notamment sur la gestion des erreurs et la sécurité multi-tenant.
Voici les 3 catégories de tests les plus importantes à ajouter :
1. Les scénarios d'erreur HTTP (Le plus urgent)
Ton handler actuel renvoie souvent une erreur 500 générique via api.RespondWithError(w, err) car il ne distingue pas encore les types d'erreurs (ton commentaire // TODO dans GetByID le confirme).
Il faut des tests qui provoquent ces erreurs pour t'obliger à implémenter le mapping des statuts HTTP corrects.
Test GetByID "Not Found" :
Action : Appeler GET /users/{id-inexistant}.
Attendu : Statut HTTP 404 Not Found. (Actuellement, ton code renverrait probablement une 500, ce test échouera donc au début, ce qui est le but).
Test Create "Bad Request" :
Action : Appeler POST /users avec un JSON invalide (ex: une virgule en trop) ou des champs obligatoires manquants (si le service les valide).
Attendu : Statut HTTP 400 Bad Request.
2. La sécurité Multi-tenant (Isolation)
Tu dois absolument vérifier qu'un locataire ne peut pas accéder aux données d'un autre. C'est la faille la plus critique dans un SaaS.
Test d'isolation inter-tenant :
Setup : Créer un utilisateur "Arthur" dans le tenant-A (dans le fake repo).
Action : Tenter de récupérer Arthur via GET /users/{arthur-id} en simulant un contexte avec le tenant-B.
Attendu : Statut HTTP 404 Not Found (Arthur n'existe pas dans le tenant B).
3. Le flux de recherche (Search)
Tu as implémenté le handler et le fake repository pour Search, mais il n'y a aucun test.
Test Search_Flow :
Setup : Peupler le fake repo avec 3 utilisateurs (2 "Arthur" et 1 "Léodagan") dans le même tenant.
Action : Appeler GET /users?nom=Arthur.
Attendu : Statut 200 OK et un tableau JSON contenant exactement les 2 utilisateurs "Arthur".
En résumé : Priorise les tests qui vérifient que ton API échoue "proprement" (404, 400, 401/403) et qu'elle ne mélange pas les données des locataires.




POUR LE REFACTORING : 


Oui, un refactoring est nécessaire pour ne pas te décourager d'écrire des tests.
Voici les deux techniques standard en Go pour alléger ça :
1. Créer une fonction "Helper de Setup"
Au lieu de répéter les 10 premières lignes dans chaque test, tu crées une fonction privée qui te renvoie tout ce qui est déjà configuré.
Conceptuellement, ça donnerait ça :
code
Go
// (Pseudo-code conceptuel, je ne génère pas le code final)
func setupTest(t *testing.T) (*fakeRepo, *chi.Mux) {
    // Instancier repo, service, handler
    // Configurer un routeur Chi de test
    // Renvoyer le repo (pour les assertions) et le routeur (pour faire les requêtes)
}

func TestCreateUser(t *testing.T) {
    repo, router := setupTest(t) // Hop, 10 lignes économisées

    // ... suite du test plus directe ...
}
2. Les "Table-Driven Tests" (Tests pilotés par des tableaux)
C'est LA méthode idiomatique en Go pour éviter la répétition quand on teste plusieurs variantes d'une même fonctionnalité (surtout pour les cas d'erreurs et de validation).
Au lieu d'avoir 5 fonctions TestCreateUser_InvalidEmail, TestCreateUser_MissingName, etc., tu fais une seule fonction qui boucle sur un tableau de scénarios.
Structure conceptuelle :
code
Go
func TestCreateUser_Scenarios(t *testing.T) {
    // On définit la structure d'un scénario de test
    tests := []struct {
        name           string      // Nom du scénario (ex: "Email invalide")
        inputBody      interface{} // Données envoyées
        tenantID       string      // Contexte
        expectedStatus int         // Statut HTTP attendu (ex: 400)
    }{
        {"Succès standard", validBody, "tenant-1", 201},
        {"Email invalide", badEmailBody, "tenant-1", 400},
        {"Tenant manquant", validBody, "", 401},
    }

    // On initialise une seule fois l'environnement
    repo, router := setupTest(t)

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            // Ici, la logique d'appel HTTP standard
            // On utilise tc.inputBody, tc.tenantID...
            // On vérifie que le statut reçu == tc.expectedStatus
        })
    }
}
Conclusion :
Adopte ces deux patterns. Tes fichiers de tests seront deux fois plus courts et beaucoup plus faciles à lire et à maintenir.
Si tu veux, je peux te montrer l'application concrète de ces deux patterns sur ton user_test.go.