.PHONY: sort generator validator

GOBIN = ./bin
GOCMD = ./cmd

define build_cmd
	@echo "Building $(1) cmd"

	go build -o $(GOBIN)/$(1) $(GOCMD)/$(1)

	@echo
	@echo "Run \"$(GOBIN)/$(1)\" to launch $(1)"
	@echo
endef

sort:
	$(call build_cmd,sort)

generator:
	$(call build_cmd,generator)

validator:
	$(call build_cmd,validator)

all: sort generator validator