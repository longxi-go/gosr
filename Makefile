.DEFAULT_GOAL := help

GO ?= go
UV ?= uv
GOLANGCI_LINT ?= golangci-lint
TOOLS_MODE ?= source
RELEASE_MODULE ?= root
REPOSITORY ?= longxi-go/gosr
TAG ?=
export GO GOLANGCI_LINT

.PHONY: help setup format check-format test-format license-headers check-license-headers
.PHONY: test race vet lint verify-mod verify-tools verify-tools-source verify-tools-published
.PHONY: test-scripts validate-release prepare-tools-release check check-all check-published compatibility
.PHONY: docs-setup docs-serve docs-build

help:
	@printf '%s\n' 'Usage: make <target>' '' \
	  '  setup                    Resolve and verify pinned golayout tooling' \
	  '  license-headers          Refresh owned Go source license headers' \
	  '  check-license-headers    Check source headers without writes' \
	  '  format                   Rewrite owned Go files with golayout' \
	  '  check-format             Check source layout without writes' \
	  '  test-format              Test shared formatting and selection gates' \
	  '  test / race / vet / lint Check both modules against this checkout' \
	  '  verify-mod               Check root module tidiness independently' \
	  '  verify-tools-source      Check tools tidiness against this checkout' \
	  '  verify-tools-published   Check tools tidiness against published modules' \
	  '  compatibility            Run module, vet and test compatibility gates' \
	  '  test-scripts             Test release and module-check scripts' \
	  '  validate-release         Validate TAG=vX.Y.Z [RELEASE_MODULE=root|tools]' \
	  '  prepare-tools-release    Resolve published dependencies for TAG=tools/vX.Y.Z' \
	  '  check                    Run all source quality gates' \
	  '  check-all                Alias for check' \
	  '  check-published          Run all gates using published tool dependencies' \
	  '  docs-setup / docs-serve / docs-build  Manage the bilingual manual'

setup:
	@python3 scripts/golayout.py setup

license-headers:
	@python3 scripts/golayout.py headers-write

check-license-headers:
	@python3 scripts/golayout.py headers-check

format:
	@python3 scripts/golayout.py write

check-format:
	@python3 scripts/golayout.py check

test-format:
	@python3 scripts/golayout.py test

test:
	GOWORK=off $(GO) test -cover ./...
	python3 scripts/tools.py $(TOOLS_MODE) -- $(GO) test -cover ./...

race:
	GOWORK=off $(GO) test -race -shuffle=on ./...
	python3 scripts/tools.py $(TOOLS_MODE) -- $(GO) test -race -shuffle=on ./...

vet:
	GOWORK=off $(GO) vet ./...
	python3 scripts/tools.py $(TOOLS_MODE) -- $(GO) vet ./...

lint:
	GOWORK=off $(GOLANGCI_LINT) run --timeout=10m
	python3 scripts/tools.py $(TOOLS_MODE) -- $(GOLANGCI_LINT) run --timeout=10m --config ../.golangci.yml

verify-mod:
	GOWORK=off $(GO) mod tidy -diff

verify-tools-source:
	python3 scripts/tools.py source -- $(GO) mod tidy -diff

verify-tools-published:
	python3 scripts/tools.py published -- $(GO) mod tidy -diff

verify-tools:
	python3 scripts/tools.py $(TOOLS_MODE) -- $(GO) mod tidy -diff

test-scripts:
	python3 -m unittest discover -s scripts/tests -v

validate-release:
	python3 scripts/release.py --module "$(RELEASE_MODULE)" --tag "$(TAG)"

prepare-tools-release:
	python3 scripts/release.py --module tools --tag "$(TAG)" --check-root-release --repository "$(REPOSITORY)"
	python3 scripts/tools.py published -- $(GO) mod tidy

compatibility: verify-mod verify-tools vet test
check: check-license-headers check-format test-format verify-mod verify-tools vet test race lint test-scripts
check-all: check

check-published:
	$(MAKE) --no-print-directory check TOOLS_MODE=published

docs-setup:
	cd docs/design && $(UV) sync --locked

docs-serve: docs-setup
	cd docs/design && $(UV) run --locked mkdocs serve

docs-build: docs-setup
	cd docs/design && $(UV) run --locked mkdocs build --strict
	cd docs/design && $(UV) run --locked python -m unittest discover -s tests -v

.DEFAULT:
	@printf 'Unknown target: %s. Available targets follow.\n' "$@" >&2
	@$(MAKE) --no-print-directory help
	@exit 2
