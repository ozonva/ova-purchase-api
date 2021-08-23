default: build


download:
	go mod download

build: download
	go build -o bin/main cmd/ova-purchase-api/main.go

generate: download
	mockgen -destination=./internal/mocks/repo_mock.go -source internal/repo/repo.go

test: download
	go test ./...

clean:
	rm -rf ./bin
	go clean -testcache
