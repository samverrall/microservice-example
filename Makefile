OUTPUT := bin/service

.PHONY: lint
lint:
	golangci-lint --enable gosec,misspell run ./...

.PHONY: run
run:
	./run.sh

.PHONY: build
build: lint
	go build -o $(OUTPUT) github.com/samverrall/microservice-example/cmd

.PHONY: test
test:
	go test -v --race ./...

.PHONY: govulncheck
govulncheck:
	govulncheck ./...

.PHONY: gogen
gogen:
	go generate ./...

.PHONY: proto-health
proto-health:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
    --go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
    pkg/proto/health.proto

.PHONY: proto-all
proto-all:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
    --go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
    pkg/proto/*.proto
