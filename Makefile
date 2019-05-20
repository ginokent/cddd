SHELL := $(shell command -v bash)

COMMAND_NAME := cddd
MAIN_PATH := ./cmd/$(COMMAND_NAME)
BUILD_PATH := ./_build

.PHONY: help
help:  ## display this documents
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: setup-tools
setup-tools:  ## setup tools for development
	# If not found in $PATH, install
	command -v goxz || (cd ~ && go get github.com/Songmu/goxz/cmd/goxz)

.PHONY: generate
generate:  ## Execute `go generate`
	# generate
	go generate -v ./...

.PHONY: run
run:  ## Execute `go run` with $(ARGS)
	##
	# You may have to run \`go generate\` manually.
	##
	# run
	go run $(MAIN_PATH) $(ARGS)

.PHONY: package
package: generate ## build executable file to $(BUILD_PATH)
	# build
	goxz -pv $(shell go run $(MAIN_PATH) --version | cut -d' ' -f3) -d $(BUILD_PATH) $(MAIN_PATH)

.PHONY: install
install: generate ## install executable file to ${GOPATH}/bin
	# install
	go install $(MAIN_PATH)

.PHONY: test
test: generate ## Execute `go test`
	# test
	go test -v -cover ./...

.PHONY: complete
complete:  ## for bash complete
	# complete
	# Append the following to bashrc for bash completion.
	@echo "test -f ~/.urfave_cli_bash_autocomplete || curl -LRsS https://raw.githubusercontent.com/urfave/cli/master/autocomplete/bash_autocomplete -o ~/.urfave_cli_bash_autocomplete"
	@echo "PROG=cddd source ~/.urfave_cli_bash_autocomplete"

