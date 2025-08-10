# Default target
.PHONY: default
default: all
# Define the Go binary name
BINARY_NAME_UBI7 = dist/kpasscli-ubi7
BINARY_NAME_UBI8 = dist/kpasscli-ubi8
BINARY_NAME = dist/kpasscli

# Build the binary

${BINARY_NAME}-linux-amd64: $(shell find . -type f -name '*.go')
	mkdir -p dist/linux-amd64
	GOOS=linux GOARCH=amd64 go build -v -o dist/linux-amd64/kpasscli

${BINARY_NAME}-windows-amd64: $(shell find . -type f -name '*.go')
	mkdir -p dist/windows-amd64
	GOOS=windows GOARCH=amd64 go build -v -o dist/windows-amd64/kpasscli.exe

${BINARY_NAME}-darwin-amd64: $(shell find . -type f -name '*.go')
	mkdir -p dist/darwin-amd64
	GOOS=darwin GOARCH=amd64 go build -v -o dist/darwin-amd64/kpasscli

${BINARY_NAME}-darwin-arm64: $(shell find . -type f -name '*.go')
	mkdir -p dist/darwin-arm64
	GOOS=darwin GOARCH=arm64 go build -v -o dist/darwin-arm64/kpasscli

build: ${BINARY_NAME}-linux-amd64 ${BINARY_NAME}-windows-amd64 ${BINARY_NAME}-darwin-amd64 ${BINARY_NAME}-darwin-arm64

build-ubi7:
	docker build -t kpasscli:ubi7 -f Dockerfile-ubi7 .
	if ! docker ps -a | rg 'kpasscli-container' >/dev/null; then \
		docker create --name kpasscli-container kpasscli:ubi7; \
	else \
		docker rm -f kpasscli-container; \
		docker create --name kpasscli-container kpasscli:ubi7; \
	fi
	docker cp kpasscli-container:/app/dist/kpasscli $(BINARY_NAME_UBI7)
# 	scp $(BINARY_NAME_UBI7) cid-scp0-tls-v01-mgmt:
# 	artifactory-upload.sh -lf=$(BINARY_NAME_UBI7) -tr=scptools-bin-dev-local  -tf=ocp-stable-4.16/clients/oc/4.16/kpasscli

# Clean up build artifacts
clean:
	rm -rf dist/*

# Clean up build artifacts
clean-ubi7:
	rm -f $(BINARY_NAME_UBI7)

# Run the application
run: build
	./$(BINARY_NAME)

# Test the application
test: build
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

coverage: test
	go tool cover -html=coverage.out

# Test the application in a UBI7 container
test-ubi7:
	docker build -t kpasscli:ubi7 -f Dockerfile-ubi7 .
	docker run --rm --entrypoint=/bin/bash kpasscli:ubi7 -c "go mod tidy && go test ./..."


# Generate documentation
doc:
	mkdir -p dist/doc
	gomarkdoc --output dist/doc/README.md ./...

# Generate a static site for documentation (not supported by gomarkdoc, so just copy doc as placeholder)
docsite: doc
	mkdir -p dist/docsite
	cp dist/doc/README.md dist/docsite/index.md


# Generate a graphical representation of the Go module dependency graph
modgraph:
	go mod graph | sed 's/@[^ ]*//g' | awk '{gsub(/[^a-zA-Z0-9._/-]/, "_", $$1); gsub(/[^a-zA-Z0-9._/-]/, "_", $$2); print "\"" $$1 "\" -> \"" $$2 "\";"}' | \
	{ echo "digraph G {"; cat; echo "}"; } | dot -Tpng -o dist/go_mod_graph.png

all: clean build test doc docsite modgraph