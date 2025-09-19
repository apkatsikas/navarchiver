bootstrap-test:
	~/go/bin/ginkgo bootstrap

bootstrap-test:
	~/go/bin/ginkgo generate

run-tests:
	go test -p 1 -coverprofile coverage.out ./...

run-single-test:
	~/go/bin/ginkgo --focus "test 1" testDir

coverage:
	@go tool cover -html coverage.out -o coverage.html
	explorer.exe coverage.html

mocks:
	go generate ./...

build-pi:
	CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOARCH=arm64 go build -o bin/navarchiver -ldflags="-extldflags=-static" -tags sqlite_omit_load_extension ./cmd

check-formatting:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files need formatting:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

install-staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest

check-build:
	go build -o navarchiver -race ./cmd/main.go
