# Makefile à la racine du projet

.PHONY: help run test lint

# Astuce : cette commande "help" parse le Makefile lui-même pour afficher la doc
help: ## Affiche cette aide
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Lance le serveur en local (nécessite les variables d'env)
	go run main.go

test: ## Lance tous les tests (unitaires + intégration)
	go test ./...

test-short: ## Lance uniquement les tests unitaires (rapide)
	go test -short ./...

lint: ## Lance le linter golangci-lint
	golangci-lint run