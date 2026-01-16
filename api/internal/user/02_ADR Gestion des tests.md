# ADR: [Gestion des tests]

**Date:** 16/01/2026

## Context
Comment tester proprement mon domaine, en sachant que je cherche comme toujours quelque chose de simple qui ne soit pas compliqué

## Decision
Un seul fichier qui test handler et service en même temps, pour éviter que les tests soient trop long à écrire.
On fake le repository; ce genre de test pourrait s'appeler Test de Composant et surtout de COMPORTEMENT.

## Consequences
**Positif :**
Feedback immédiat pendant le développement (tests très rapides).
Réduction du code de test (on teste la fonctionnalité de bout en bout du domaine).
Aucune dépendance externe requise pour lancer les tests métier.

**Négatif :**
Les requêtes réelles vers la base de données (SQL ou Cosmos) ne sont pas testées. Il existe un risque que le code fonctionne avec le Fake mais échoue avec la vraie DB (ex: erreur de syntaxe, contrainte d'intégrité).
La maintenance des Fakes peut devenir complexe si la logique de recherche devient avancée.

## Alternatives Considered
Simuler une cosmos avec un tests container, mais ça ne fonctionnais pas correctement, et au final trop lourd des tests qui font tout.

