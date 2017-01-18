TRIS := $(CURDIR)/../../cloudflare/tls-tris

.PHONY: bin/crypto-tls-bogo-shim
bin/crypto-tls-bogo-shim:
	GOBIN=$(CURDIR)/bin $(TRIS)/_dev/go.sh install -v .

.PHONY: run
run: bin/crypto-tls-bogo-shim
ifdef NAME
	$(eval ARGS += -test $(NAME))
endif
ifdef EXTRA_ARGS
	$(eval ARGS += $(EXTRA_ARGS))
endif
	unset GOPATH && cd vendor/github.com/google/boringssl/ssl/test/runner && \
	go test -loose-errors -allow-unimplemented -shim-path $(CURDIR)/bin/crypto-tls-bogo-shim -shim-config $(CURDIR)/config.json $(ARGS)
