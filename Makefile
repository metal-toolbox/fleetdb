all: lint test
PHONY: test coverage lint golint clean vendor local-dev-databases docker-up docker-down integration-test unit-test
GOOS=linux
DB_STRING=host=localhost port=26257 user=root sslmode=disable
DEV_DB=${DB_STRING} dbname=fleetdb
TEST_DB=${DB_STRING} dbname=fleetdb_test
DOCKER_IMAGE := "ghcr.io/metal-toolbox/fleetdb"
PROJECT_NAME := fleetdb
REPO := "https://github.com/metal-toolbox/fleetdb.git"
SQLBOILER := v4.15.0
SQLBOILER_DRIVER := v4.0.0

## run all tests
test: | unit-test integration-test

## run integration tests
integration-test: test-database
	@echo Running integration tests...
	@FLEETDB_CRDB_URI="${TEST_DB}" go test -cover -tags testtools,integration -p 1 -timeout 2m ./... | \
	grep -v "could not be registered in Prometheus\" error=\"duplicate metrics collector registration attempted\"" # TODO; Figure out why this message spams when tests fail

## run unit tests
unit-test: | test-database
	@echo Running unit tests...
	@FLEETDB_CRDB_URI="${TEST_DB}" go test -cover -short -tags testtools ./...

## run single integration test. Example: make single-test test=TestIntegrationServerListComponents
single-test:
	@FLEETDB_CRDB_URI="${TEST_DB}" go test -timeout 30s -tags testtools -run ^${test}$$ github.com/metal-toolbox/fleetdb/pkg/api/v1 -v

## check test coverage
coverage: | test-database
	@echo Generating coverage report...
	@FLEETDB_CRDB_URI="${TEST_DB}" go test ./... -race -coverprofile=coverage.out -covermode=atomic -tags testtools,integration -p 1
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out

## lint
lint: | vendor
	@echo Linting Go files...
	@golangci-lint run

## clean docker files
clean: docker-clean test-clean
	@echo Cleaning...
	@rm -rf ./dist/
	@rm -rf coverage.out

## clean test env
test-clean:
	@go clean -testcache

## download/tidy go modules
vendor:
	@go mod download
	@go mod tidy

## setup docker compose test env
docker-up:
	@docker-compose -f quickstart.yml up -d crdb

## stop docker compose test env
docker-down:
	@docker-compose -f quickstart.yml down

## clean docker volumes
docker-clean:
	@docker-compose -f quickstart.yml down --volumes

## setup devel database
dev-database: | vendor
	@cockroach sql --insecure -e "drop database if exists fleetdb"
	@cockroach sql --insecure -e "create database fleetdb"
	@FLEETDB_CRDB_URI="${DEV_DB}" go run main.go migrate up

## setup test database
test-database: | vendor docker-up
	@cockroach sql --insecure -e "drop database if exists fleetdb_test"
	@cockroach sql --insecure -e "create database fleetdb_test"
	@FLEETDB_CRDB_URI="${TEST_DB}" go run main.go migrate up
	@cockroach sql --insecure -e "use fleetdb_test; ALTER TABLE attributes DROP CONSTRAINT check_server_id_server_component_id; ALTER TABLE versioned_attributes DROP CONSTRAINT check_server_id_server_component_id;"

## purge dev environment, build new image, and run tests
fresh-test: clean
	@make push-image-devel
	@make docker-up
	@make test

## install sqlboiler
install-sqlboiler:
	go install github.com/volatiletech/sqlboiler/v4@${SQLBOILER}
	go install github.com/metal-toolbox/sqlboiler-crdb-fleetdb/v4@${SQLBOILER_DRIVER}

## boil sql, if you get this error (server closed the connection), try again.
boil: install-sqlboiler test-database
	sqlboiler crdb-fleetdb --add-soft-deletes

## log into database
psql:
	@psql -d "${TEST_DB}"

## Build linux bin
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${PROJECT_NAME}

## build docker image and tag as ghcr.io/metal-toolbox/fleetdb:latest
build-image: build-linux
	docker build --rm=true -f Dockerfile -t ${DOCKER_IMAGE}:latest . \
		--label org.label-schema.schema-version=1.0 \
		--label org.label-schema.vcs-ref=${GIT_COMMIT_FULL} \
		--label org.label-schema.vcs-url=${REPO}

## build and push devel docker image to KIND image repo used by the sandbox - https://github.com/metal-toolbox/sandbox
push-image-devel: build-image
	docker tag ${DOCKER_IMAGE}:latest localhost:5001/${PROJECT_NAME}:latest
	docker push localhost:5001/${PROJECT_NAME}:latest
	kind load docker-image localhost:5001/${PROJECT_NAME}:latest

# https://gist.github.com/prwhite/8168133
# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=20

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-${TARGET_MAX_CHAR_NUM}s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' ${MAKEFILE_LIST}
