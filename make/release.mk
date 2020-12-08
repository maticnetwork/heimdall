BUILDDIR ?= $(CURDIR)/build

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

heimdall:
	$(GO) build $(BUILD_FLAGS) -o ./build/heimdalld ./cmd/heimdalld

.PHONY: build heimdall