include .envrc
# ----Development----
## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api --port=8080 --db-dsn=${BASEDWEB_DB_DSN}