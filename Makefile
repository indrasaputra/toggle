GO_UNIT_TEST_FILES	= $(shell go list ./... | grep -v /feature)
PROTOGEN_IMAGE 		= indrasaputra/protogen:2021-12-10

include Makefile.help.mk

##@ format
.PHONY: tidy
tidy: ## Tidy up go module.
	GO111MODULE=on go mod tidy

.PHONY: format
format: ## Format golang and proto files.
	bin/format.sh

.PHONY: lint.cleancache
lint.cleancache: ## Clean golangci-lint cache.
	golangci-lint cache clean

.PHONY: lint
lint: ## Lint proto files using buf and golang files using golangci-lint.
lint: lint.cleancache
	buf lint
	protolint lint -fix .
	golangci-lint run ./...

.PHONY: pretty
pretty: ## Prettify golang and proto files. Basically, it runs tidy, format, and lint command.
pretty: tidy format lint

.PHONY: check.import
check.import: ## Check if import blocks are separated accordingly.
	bin/check-import.sh

.PHONY: precommit
precommit: ## Run pre-commit to all files.
	pre-commit run --all-files

##@ Generator
.PHONY: gen.proto
gen.proto: ## Generate golang files from proto.
	bin/generate-proto.sh

.PHONY: gen.proto.docker
gen.proto.docker: ## Generate golang files from proto using docker.
	docker run -it --rm \
    --mount "type=bind,source=$(PWD),destination=/work" \
    --mount "type=volume,source=toggle-go-mod-cache,destination=/go,consistency=cached" \
    --mount "type=volume,source=toggle-buf-cache,destination=/home/.cache,consistency=cached" \
    -w /work $(PROTOGEN_IMAGE) make -e -f Makefile gen.proto pretty

.PHONY: gen.mock
gen.mock: ## Generate mock from all golang interfaces.
	bin/generate-mock.sh

##@ Build
.PHONY: compile
compile: ## Compile golang code to binary.
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o toggle cmd/server/main.go

.PHONY: build.envoy
build.envoy: ## Build docker envoy.
	docker build --no-cache -t envoy/toggle:latest -f envoy.dockerfile .

##@ Test
.PHONY: test.unit
test.unit: ## Run unit test.
test.unit: test.cleancache
	go test -failfast -v -race $(GO_UNIT_TEST_FILES)

.PHONY: test.cover
test.cover: ## Run unit test with coverage status printed in stdout.
test.cover: test.cleancache
	go test -failfast -v -race $(GO_UNIT_TEST_FILES) -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func coverage.out

.PHONY: test.coverhtml
test.coverhtml: ## Run unit test with coverage status outputed as HTML.
test.coverhtml: test.cleancache
	go test -failfast -v -race $(GO_UNIT_TEST_FILES) -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: test.cleancache
test.cleancache: ## Clean test cache.
	go clean -testcache

.PHONY: test.integration
test.integration: ## Run end-to-end test using Godog.
	bin/godog.sh

.PHONY: test.smoke
test.smoke: ## Run smoke test using k6.
	docker run -it --rm \
	--mount "type=bind,source=$(PWD),destination=/work" \
	-w /work loadimpact/k6:0.33.0 run internal/script/loadtest/smoke_test.js

.PHONY: test.load
test.load: ## Run load test using k6.
	docker run -it --rm \
	--mount "type=bind,source=$(PWD),destination=/work" \
	-w /work loadimpact/k6:0.33.0 run internal/script/loadtest/load_test.js

.PHONY: test.stress
test.stress: ## Run stress test using k6.
	docker run -it --rm \
	--mount "type=bind,source=$(PWD),destination=/work" \
	-w /work loadimpact/k6:0.33.0 run internal/script/loadtest/stress_test.js

##@ Migration
.PHONY: migration
migration: ## Create database migration.
	migrate create -ext sql -dir db/migrations $(name)

.PHONY: migrate
migrate: ## Run database migrations.
	migrate -path db/migrations -database "$(url)" -verbose up

.PHONY: rollback
rollback: ## Rollback one migration.
	migrate -path db/migrations -database "$(url)" -verbose down 1

.PHONY: rollback.all
rollback.all: ## Rollback all migrations.
	migrate -path db/migrations -database "$(url)" -verbose down -all

.PHONY: migrate.force
migrate.force: ## Force migrate specific version.
	migrate -path db/migrations -database "$(url)" -verbose force $(version)

.PHONY: validate.migration
validate.migration: ## Validate migration files.
	bin/validate-migration.sh
