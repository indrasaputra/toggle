GO_UNIT_TEST_FILES	= $(shell go list ./... | grep -v /feature)
PROTOGEN_IMAGE 		= indrasaputra/protogen:0.0.1

.PHONY: tidy
tidy:
	GO111MODULE=on go mod tidy

.PHONY: format
format:
	bin/format.sh

.PHONY: lint.cleancache
lint.cleancache:
	golangci-lint cache clean

.PHONY: lint
lint: lint.cleancache
	buf lint
	golangci-lint run ./...

.PHONY: pretty
pretty: tidy format lint

.PHONY: gen.proto
gen.proto:
	bin/generate-proto.sh

.PHONY: gen.proto.docker
gen.proto.docker:
	docker run -it --rm \
    --mount "type=bind,source=$(PWD),destination=/work" \
    --mount "type=volume,source=toggle-go-mod-cache,destination=/go,consistency=cached" \
    --mount "type=volume,source=toggle-buf-cache,destination=/home/.cache,consistency=cached" \
    -w /work $(PROTOGEN_IMAGE) make -e -f Makefile gen.proto pretty

.PHONY: gen.mock
gen.mock:
	bin/generate-mock.sh

.PHONY: check.import
check.import:
	bin/check-import.sh

.PHONY: compile
compile:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o toggle cmd/server/main.go

.PHONY: test.unit
test.unit: test.cleancache
	go test -v -race $(GO_UNIT_TEST_FILES)

.PHONY: test.cover
test.cover:
	go test -v -race $(GO_UNIT_TEST_FILES) -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func coverage.out

.PHONY: test.coverhtml
test.coverhtml:
	go test -v -race $(GO_UNIT_TEST_FILES) -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: test.cleancache
test.cleancache:
	go clean -testcache

.PHONY: test.integration
test.integration:
	bin/godog.sh

.PHONY: test.smoke
test.smoke:
	docker run -it --rm \
	--mount "type=bind,source=$(PWD),destination=/work" \
	-w /work loadimpact/k6:0.33.0 run internal/script/loadtest/smoke_test.js

.PHONY: test.load
test.load:
	docker run -it --rm \
	--mount "type=bind,source=$(PWD),destination=/work" \
	-w /work loadimpact/k6:0.33.0 run internal/script/loadtest/load_test.js

.PHONY: test.stress
test.stress:
	docker run -it --rm \
	--mount "type=bind,source=$(PWD),destination=/work" \
	-w /work loadimpact/k6:0.33.0 run internal/script/loadtest/stress_test.js

.PHONY: migration
migration:
	migrate create -ext sql -dir db/migrations $(name)

.PHONY: migrate
migrate:
	migrate -path db/migrations -database "$(url)?sslmode=disable" -verbose up

.PHONY: rollback
rollback:
	migrate -path db/migrations -database "$(url)?sslmode=disable" -verbose down 1

.PHONY: rollback.all
rollback.all:
	migrate -path db/migrations -database "$(url)?sslmode=disable" -verbose down -all

.PHONY: migrate.force
migrate.force:
	migrate -path db/migrations -database "$(url)?sslmode=disable" -verbose force $(version)

.PHONY: validate.migration
validate.migration:
	bin/validate-migration.sh
