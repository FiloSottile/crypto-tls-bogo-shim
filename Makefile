bin/crypto-tls-bogo-shim: *.go
	GOBIN=$(CURDIR)/bin go1.8beta2 install .

.PHONY: run
run: bin/crypto-tls-bogo-shim
	unset GOPATH && cd vendor/github.com/google/boringssl/ssl/test/runner && \
	go test -loose-errors -allow-unimplemented -shim-path $(CURDIR)/bin/crypto-tls-bogo-shim -shim-config $(CURDIR)/config.json
