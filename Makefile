KEY_FILE = certs/key.pem
CERT_FILE = certs/cert.pem
CONFIG_FILE = testdata/cert.cnf

all: build

build:
	go build -v ./...

test: generate-keys
	go test -race -coverprofile=coverage.txt -covermode=atomic

generate-keys:
	# Generate the TLS certificate and private key.
	mkdir -p certs
	openssl req -new -x509 -sha256 -days 365 -nodes -keyout $(KEY_FILE) -out $(CERT_FILE) -config $(CONFIG_FILE)

.PHONY: all build test generate-keys
