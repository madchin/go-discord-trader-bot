build-dev:
	docker compose run --rm --name DEBUG -e RUNTIME_ENVIRONMENT=DEV -d trader
build-prod:
	docker compose up -d