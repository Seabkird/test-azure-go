# ADR: [Choix cosmos]

**Date:** 15/01/2026

## Context
Comment persister les données 

## Decision
Utiliser du no sql et nottament cosmosDB, je part du principe que le code go est source de vérité absolu. 
De plus la cosmosDB existe à la consommation. 
Le noSQL permet de ne pas avoir à gérer le db est permet juste de refléter le code écrit en go. 

## Consequences
Plus compliqué de passer de BASE à ACID

## Alternatives Considered
Faire du postgre sql, mais ça m'a semblé  trop compliqué  et couteux.