# Define the Go binary name
BINARY_NAME_UBI7 = dist/kpasscli-ubi7
BINARY_NAME_UBI8 = dist/kpasscli-ubi8
BINARY_NAME = dist/kpasscli

# Build the binary
build:
	go mod tidy
	go build -v -o $(BINARY_NAME)

build-ubi7:
	podman build -t kpasscli:ubi7 .
	if ! podman ps -a | rg 'kpasscli-container' >/dev/null; then \
		podman create --name kpasscli-container kpasscli:ubi7; \
	fi
	podman cp kpasscli-container:/app/dist/kpasscli $(BINARY_NAME_UBI7)
	scp $(BINARY_NAME_UBI7) cid-scp0-tls-v01-mgmt:
	artifactory-upload.sh -lf=$(BINARY_NAME_UBI7) -tr=scptools-bin-dev-local  -tf=ocp-stable-4.16/clients/oc/4.16/kpasscli

# Clean up build artifacts
clean:
	rm -f $(BINARY_NAME)

# Clean up build artifacts
clean-ubi7:
	rm -f $(BINARY_NAME_UBI7)

# Run the application
run: build
	./$(BINARY_NAME)

# Generate a graphical representation of the Go module dependency graph
modgraph:
	go mod graph | sed 's/@[^ ]*//g' | awk '{gsub(/[^a-zA-Z0-9._/-]/, "_", $$1); gsub(/[^a-zA-Z0-9._/-]/, "_", $$2); print "\"" $$1 "\" -> \"" $$2 "\";"}' | \
	{ echo "digraph G {"; cat; echo "}"; } | dot -Tpng -o dist/go_mod_graph.png
