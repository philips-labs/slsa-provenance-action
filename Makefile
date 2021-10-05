PKG=github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli

GO_BUILD_FLAGS := -trimpath
COMMANDS       := slsa-provenance

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

FORCE: ;

bin/%: cmd/% FORCE
	CGO_ENABLED=0 go build $(GO_BUILD_FLAGS) -o $@ ./$<

.PHONY: download
download: ## download dependencies via go mod
	go mod download

.PHONY: build
build: $(addprefix bin/,$(COMMANDS)) ## builds binaries

.PHONY: image
image: ## build the binary in a docker image
	docker build \
		-t "philipssoftware/slsa-provenance:$(GIT_TAG)" \
		-t "philipssoftware/slsa-provenance:$(GIT_HASH)" .
