# TODO REST API User

- [ ] **Refactorisation HTTP Helpers** : Déplacer `respondWithJSON`, `respondWithError` et `getTenantIDFromContext` dans un package partagé (ex: `internal/kit/http`).
- [ ] **Erreurs Structurées** : Créer une struct `APIError` standard (code, message) au lieu d'une simple string.
- [ ] **Sécurité Erreurs 500** : Ne jamais renvoyer l'erreur brute `err.Error()` en prod. Logger l'erreur et renvoyer un message générique au client.
- [ ] **Mapping Erreurs Métier** : Dans les handlers, vérifier le type d'erreur retourné par le service (ex: `ErrUserAlreadyExists`, `ErrNotFound`) pour renvoyer le bon status code (409, 404) au lieu de 500 systématiquement.
- [ ] **Améliorer Search Parsing** : Automatiser le parsing des paramètres d'URL dans `parseSearchFilter` (envisager la lib `github.com/go-playground/form`) pour éviter le code répétitif.