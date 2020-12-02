.PHONY: bins
bins: $(BINS)

.PHONY: build
build:
	$(GO) build -o ./build ./cmd/heimdalld

heimdall:
	$(GO) build $(BUILD_FLAGS) -o ./build ./cmd/heimdalld

install:
	$(GO) install $(BUILD_FLAGS) ./cmd/heimdalld
