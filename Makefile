# Default target
.PHONY: default
default: build
# Define the Go binary name
BINARY_NAME_UBI7 = dist/kpasscli-ubi7
BINARY_NAME_UBI8 = dist/kpasscli-ubi8
BINARY_NAME = dist/kpasscli

# Build the binary

$(BINARY_NAME): $(shell find . -type f -name '*.go')
	go mod tidy
	go build -v -o $(BINARY_NAME) .

build: $(BINARY_NAME)

build-ubi7:
	docker build -t kpasscli:ubi7 -f Dockerfile-ubi7 .
	if ! docker ps -a | rg 'kpasscli-container' >/dev/null; then \
		docker create --name kpasscli-container kpasscli:ubi7; \
	else \
		docker rm -f kpasscli-container; \
		docker create --name kpasscli-container kpasscli:ubi7; \
	fi
# 	docker cp kpasscli-container:/app/dist/kpasscli $(BINARY_NAME_UBI7)
# 	scp $(BINARY_NAME_UBI7) cid-scp0-tls-v01-mgmt:
# 	artifactory-upload.sh -lf=$(BINARY_NAME_UBI7) -tr=scptools-bin-dev-local  -tf=ocp-stable-4.16/clients/oc/4.16/kpasscli

# Clean up build artifacts
clean:
	rm -f $(BINARY_NAME)

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

# Test the application in a UBI7 container
test-ubi7: build-ubi7
	-docker rm -f kpasscli-container || true
	docker create --name kpasscli-container --entrypoint /bin/bash kpasscli:ubi7 -c "sleep infinity"
	docker start kpasscli-container
	docker exec kpasscli-container /bin/bash -c "cd /app && GO111MODULE=on go test ./..."
	docker stop kpasscli-container
	docker rm kpasscli-container

# Generate documentation
doc:
	docgen -o dist/doc -f markdown -t kpasscli

# Generate a static site for documentation
docsite:
	docgen -o dist/docsite -f static -t test_kpasscli


# Generate a graphical representation of the Go module dependency graph
modgraph:
	go mod graph | sed 's/@[^ ]*//g' | awk '{gsub(/[^a-zA-Z0-9._/-]/, "_", $$1); gsub(/[^a-zA-Z0-9._/-]/, "_", $$2); print "\"" $$1 "\" -> \"" $$2 "\";"}' | \
	{ echo "digraph G {"; cat; echo "}"; } | dot -Tpng -o dist/go_mod_graph.png
