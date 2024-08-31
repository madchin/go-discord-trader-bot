quick-setup:
	chmod +x quick_env_setup && ./quick_env_setup
build-debug:
	docker compose -f compose.dev.yaml up -d
build-prod:
	docker compose -f compose.prod.yaml up -d