GO := $(CURDIR)/vendor/github.com/cloudflare/tls-tris/_dev/go.sh

bin/crypto-tls-bogo-shim: *.go
	GOBIN=$(CURDIR)/bin $(GO) install .

.PHONY: run
run: bin/crypto-tls-bogo-shim
	unset GOPATH && cd vendor/github.com/google/boringssl/ssl/test/runner && \
	go test -loose-errors -allow-unimplemented -shim-path $(CURDIR)/bin/crypto-tls-bogo-shim -shim-config $(CURDIR)/config.json
