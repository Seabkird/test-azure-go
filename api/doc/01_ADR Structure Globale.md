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
