.PHONY: generate
generate:
	$(GO) generate ./...

.PHONY: proto-gen
proto-gen: $(PROTOC) $(GRPC_GATEWAY) $(PROTOC_GEN_COSMOS) proto-dep-setup
	bash ./scripts/protoc-gen.sh

.PHONY: proto-swagger-gen
proto-swagger-gen: $(PROTOC) $(PROTOC_SWAGGER_GEN) $(SWAGGER_COMBINE) proto-dep-setup
	bash ./scripts/protoc-swagger-gen.sh

.PHONY: update-swagger-docs
update-swagger-docs: $(STATIK) proto-swagger-gen
	$(STATIK) -src=client/docs/swagger-ui -dest=client/docs -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
		echo "\033[91mSwagger docs are out of sync!!!\033[0m"; \
		exit 1; \
	else \
		echo "\033[92mSwagger docs are in sync\033[0m"; \
	fi

.PHONY: codegen
codegen: generate proto-gen update-swagger-docs
