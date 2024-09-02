quick-setup:
	chmod +x quick_env_setup && ./quick_env_setup
lint:
	golangci-lint run
build-debug:
	docker compose -f compose.dev.yaml watch
build-prod:
	docker compose -f compose.prod.yaml up -d