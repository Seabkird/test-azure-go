TODO a alimenter proprement au fil du temps. 
Dans un ERP, les mots ont un sens très précis. Qu'est-ce qu'un "Client" ? Est-ce différent d'un "Utilisateur" ? Qu'est-ce qu'une "Commande draft" ?
Le DDD insiste sur le "Langage Ubiquitaire" (Ubiquitous Language). Le code doit utiliser ces termes.
La pratique : Créez un fichier GLOSSARY.md à la racine du projet (ou dans docs/).
Contenu : Définissez les termes métier clés.

# Glossaire Métier

| Terme | Définition dans notre ERP |
| :--- | :--- |
| **Utilisateur (User)** | Personne physique ayant un accès (login/pass) au SaaS. Défini dans le module `internal/user`. |
| **Tenant (Compte)** | L'entreprise cliente qui paie l'abonnement. Un Tenant contient plusieurs Utilisateurs. |
| **Article (Item)** | Un bien physique ou un service générique pouvant être vendu. |
| **Produit (Product)** | Un Article configuré spécifiquement pour un Tenant avec un prix de vente défini. |