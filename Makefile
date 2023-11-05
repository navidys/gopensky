BIN := ./bin
GO := go
TARGET := gopensky-query
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
PRE_COMMIT = $(shell command -v bin/venv/bin/pre-commit ~/.local/bin/pre-commit pre-commit | head -n1)
PKG_MANAGER ?= $(shell command -v dnf yum|head -n1)
BUILDFLAGS := -mod=vendor $(BUILDFLAGS)

#=================================================
# Build binary, documents
#=================================================

all: binary

.PHONY: clean
clean:
	@rm -rf $(BIN)

.PHONY: binary
binary: $(TARGET)  ## Build prometheus-podman-exporter binary
	@true

.PHONY: $(TARGET)
$(TARGET): $(SRC)
	@echo "running go build"
	@mkdir -p $(BIN)
	$(GO) build $(BUILDFLAGS) -o $(BIN)/$(TARGET) ./cmd/$(TARGET)/

.PHONY: docs
docs: ## Generates html documents
	@make -C docs


#=================================================
# Required tools installation tartgets
#=================================================

.PHONY: install.tools
install.tools: .install.pre-commit .install.codespell .install.golangci-lint ## Install needed tools

.PHONY: .install.codespell
.install.codespell:
	sudo ${PKG_MANAGER} -y install codespell

.PHONY: .install.pre-commit
.install.pre-commit:
	if [ -z "$(PRE_COMMIT)" ]; then \
		python3 -m pip install --user pre-commit; \
	fi

.PHONY: .install.golangci-lint
.install.golangci-lint:
	VERSION=1.51.1 ./hack/install_golangci.sh

#=================================================
# Linting/Formatting/Code Validation targets
#=================================================

.PHONY: validate
validate: gofmt lint pre-commit codespell vendor ## Validate prometheus-podman-exporter code (fmt, lint, ...)

.PHONY: vendor
vendor: ## Check vendor
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify
	@bash ./hack/tree_status.sh

.PHONY: lint
lint: ## Run golangci-lint
	@echo "running golangci-lint"
	$(BIN)/golangci-lint version
	$(BIN)/golangci-lint run

.PHONY: pre-commit
pre-commit:   ## Run pre-commit
ifeq ($(PRE_COMMIT),)
	@echo "FATAL: pre-commit was not found, make .install.pre-commit to installing it." >&2
	@exit 2
endif
	$(PRE_COMMIT) run -a

.PHONY: gofmt
gofmt:   ## Run gofmt
	@echo -e "gofmt check and fix"
	@gofmt -w $(SRC)

.PHONY: codespell
codespell: ## Run codespell
	@echo "running codespell"
	@codespell -S ./vendor,go.mod,go.sum,./.git,./docs/_build

#=================================================
# Help menu
#=================================================

_HLP_TGTS_RX = '^[[:print:]]+:.*?\#\# .*$$'
_HLP_TGTS_CMD = grep -E $(_HLP_TGTS_RX) $(MAKEFILE_LIST)
_HLP_TGTS_LEN = $(shell $(_HLP_TGTS_CMD) | cut -d : -f 1 | wc -L)
_HLPFMT = "%-$(_HLP_TGTS_LEN)s %s\n"
.PHONY: help
help: ## Print listing of key targets with their descriptions
	@printf $(_HLPFMT) "Target:" "Description:"
	@printf $(_HLPFMT) "--------------" "--------------------"
	@$(_HLP_TGTS_CMD) | sort | \
		awk 'BEGIN {FS = ":(.*)?## "}; \
			{printf $(_HLPFMT), $$1, $$2}'
