include .envrc
#===================#
#      Helpers      #
#===================#
## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

#=======================#
#      Development      #
#=======================#
## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api --port=8080 --db-dsn=${BASEDWEB_DB_DSN} --db-max-open-conns=30 --db-max-idle-conns=30 --db-max-idle-time="17m"
## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}
## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}
## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path=./migrations -database=${BASEDWEB_DB_DSN} up
	

# =========================#
#  	Quality Control	   #
# =========================#
## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...
## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy -compat=1.17
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor
