PROTOGEN_IMAGE = indrasaputra/protogen:v0.0.1

.PHONY: format
format:
	bin/format.sh

.PHONY: gengrpc
gengrpc:
	bin/generate-grpc.sh

.PHONY: gengrpcdocker
gengrpcdocker:
	docker run -it --rm \
    --mount "type=bind,source=$(PWD),destination=/work" \
    --mount "type=volume,source=toggle-go-mod-cache,destination=/go,consistency=cached" \
    --mount "type=volume,source=toggle-buf-cache,destination=/home/.cache,consistency=cached" \
    -w /work $(PROTOGEN_IMAGE) make -e -f Makefile gengrpc pretty

.PHONY: check-import
check-import:
	bin/check-import.sh

.PHONY: mockgen
mockgen:
	bin/generate-mock.sh

.PHONY: cleanlintcache
cleanlintcache:
	golangci-lint cache clean

.PHONY: lint
lint: cleanlintcache
	buf lint
	golangci-lint run ./...

.PHONY: tidy
tidy:
	GO111MODULE=on go mod tidy

.PHONY: pretty
pretty: tidy format lint

.PHONY: cover
cover:
	go test -v -race ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func coverage.out

.PHONY: coverhtml
coverhtml:
	go test -v -race ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: cleantestcache
cleantestcache:
	go clean -testcache

.PHONY: test-unit
test.unit: cleantestcache
	go test -v -race ./...