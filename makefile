ifneq (,$(wildcard ./.env))
    include .env
    export
endif

default: build

LOCAL_BIN:=$(CURDIR)/bin

.PHONY: run
run:
	GOBIN=$(LOCAL_BIN) go run cmd/ova-purchase-api/main.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	GOBIN=$(LOCAL_BIN) go test -v ./...

.PHONY: .build
.build:
	GOBIN=$(LOCAL_BIN) go build -o $(LOCAL_BIN)/ova-purchase-api ./cmd/ova-purchase-api/main.go


.PHONY: deps
deps: install-go-deps

.PHONY: install-go-deps
install-go-deps: .install-go-deps

.PHONY: .install-go-deps
.install-go-deps:
	ls go.mod || go mod init github.com/ozonva/ova-purchase-api
	GOBIN=$(LOCAL_BIN) go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	GOBIN=$(LOCAL_BIN) go get -u github.com/golang/protobuf/proto
	GOBIN=$(LOCAL_BIN) go get -u github.com/golang/protobuf/protoc-gen-go
	GOBIN=$(LOCAL_BIN) go get -u google.golang.org/grpc
	GOBIN=$(LOCAL_BIN) go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	GOBIN=$(LOCAL_BIN) go get -u github.com/envoyproxy/protoc-gen-validate
	GOBIN=$(LOCAL_BIN) go get -u github.com/joho/godotenv
	GOBIN=$(LOCAL_BIN) go get -u github.com/pressly/goose/v3/cmd/goose
	GOBIN=$(LOCAL_BIN) go install -u github.com/alvaroloes/enumer
	GOBIN=$(LOCAL_BIN) go install -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	GOBIN=$(LOCAL_BIN) go install -u github.com/envoyproxy/protoc-gen-validate
	GOBIN=$(LOCAL_BIN) go install -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger


.PHONY: .generate
.generate:
	mkdir -p swagger
	mkdir -p pkg/ova-purchase-api
	GOBIN=$(LOCAL_BIN) PATH=$(PATH):$(LOCAL_BIN) protoc -I vendor.protogen \
		--go_out=pkg/ova-purchase-api --go_opt=paths=import --validate_out=lang=go:pkg/ova-purchase-api \
		--go-grpc_out=pkg/ova-purchase-api --go-grpc_opt=paths=import \
		--grpc-gateway_out=pkg/ova-purchase-api \
		--grpc-gateway_opt=logtostderr=true \
		--grpc-gateway_opt=paths=import \
		--swagger_out=allow_merge=true,merge_file_name=api:swagger \
	api/ova-purchase-api/ova-purchase-api.proto
	mv pkg/ova-purchase-api/github.com/ozonva/ova-purchase-api/pkg/ova-purchase-api/* pkg/ova-purchase-api/
	rm -rf pkg/ova-purchase-api/github.com
	mkdir -p cmd/ova-purchase-api
	GOBIN=$(LOCAL_BIN) PATH=$(PATH):$(LOCAL_BIN) enumer -type=Status -json -sql internal/purchase/purchase.go
	#cd pkg/ova-purchase-api && ls go.mod || go mod init github.com/ozonva/ova-purchase-api/pkg/ova-purchase-api && go mod tidy

.PHONY: generate
generate: .vendor-proto .generate

.PHONY: vendor-proto
vendor-proto: .vendor-proto

.PHONY: .vendor-proto
.vendor-proto:
	mkdir -p vendor.protogen
	mkdir -p vendor.protogen/api/ova-purchase-api
	mkdir -p vendor.protogen/google/api
	mkdir -p vendor.protogen/envoyproxy/validate
	cp api/ova-purchase-api/ova-purchase-api.proto vendor.protogen/api/ova-purchase-api/ova-purchase-api.proto
	curl -sSL0 https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto -o vendor.protogen/google/api/http.proto
	curl -sSL0 https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto -o vendor.protogen/google/api/annotations.proto
	curl -sSL0 https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/v0.6.1/validate/validate.proto -o vendor.protogen/envoyproxy/validate/validate.proto

.PHONY: .mocks
mocks:
	mockgen -destination=./internal/mocks/repo_mock.go -source internal/repo/repo.go -package=mocks
	mockgen -destination=./internal/mocks/flusher_mock.go -source internal/flusher/flusher.go -package=mocks

.PHONY: build
build: vendor-proto .generate .build

.PHONY: migrate
migrate:
	PATH=$(PATH):$(LOCAL_BIN) goose postgres "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}" up -dir db/migration
