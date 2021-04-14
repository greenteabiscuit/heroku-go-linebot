.PHONY: up_migrate_prod
up_migrate_prod:
	API_REVISION=release go run migration/main.go up