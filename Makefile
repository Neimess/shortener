
doc:
	swag init --generalInfo cmd/shortener/main.go --output docs --parseInternal --parseDependency

run:
	bash scripts/run.sh

migration:
	bash scripts/migrate-up.sh