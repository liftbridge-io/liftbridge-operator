SHELL := /bin/bash
ROOT := $$(git rev-parse --show-toplevel)

.PHONY: codegen
codegen: vendor
	$(ROOT)/hack/codegen/codegen.sh

.PHONY: vendor
vendor:
	cd $(ROOT) && go mod vendor
