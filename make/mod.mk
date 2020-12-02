# Golang modules and vendoring

.PHONY: deps-install
deps-install:
	$(GO) mod

.PHONY: deps-tidy
deps-tidy:
	$(GO) mod tidy

.PHONY: deps-vendor
deps-vendor:
	go mod vendor
