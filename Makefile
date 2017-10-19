# colour vars
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

DOCKER_CMD=docker
GO_CMD=go
DEP_CMD=dep

GO_BUILD=$(GO_CMD) build
GO_BUILD_RACE=$(GO_CMD) build -race
GO_TEST=$(GO_CMD) test
GO_TEST_VERBOSE=$(GO_CMD) test -v
GO_INSTALL=$(GO_CMD) install -v
GO_CLEAN=$(GO_CMD) clean
GO_DEPS=$(DEP_CMD) ensure
GO_DEPS_UPDATE= $(DEP_CMD) ensure --update
GO_VET=$(GO_CMD) vet 
GO_FMT=$(GO_CMD) fmt
GO_LINT=golint

#  Packages
APP_NAME := kinesis-connectors
TOP_PACKAGE_DIR := github.com/telenor-digital-asia
PACKAGE_LIST := $(APP_NAME) \
$(APP_NAME)/emitter

build: vet
	@for p in $(PACKAGE_LIST); do \
		echo "$(OK_COLOR)==> Build $$p ...$(NO_COLOR)"; \
		$(GO_BUILD) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

clean:
	@for p in $(PACKAGE_LIST); do \
		echo "$(OK_COLOR)==> Clean $$p ...$(NO_COLOR)"; \
		$(GO_CLEAN) $(TOP_PACKAGE_DIR)/$$p; \
	done

deps:
	@for p in $(PACKAGE_LIST); do \
		echo "$(OK_COLOR)==> Install dependencies for $$p ...$(NO_COLOR)"; \
		$(GO_DEPS) || exit 1; \
	done

fmt:
	@for p in $(PACKAGE_LIST); do \
		echo "$(OK_COLOR)==> Formatting $$p ...$(NO_COLOR)"; \
		$(GO_FMT) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

lint:
	@for p in $(PACKAGE_LIST); do \
		echo "$(OK_COLOR)==> Lint $$p ...$(NO_COLOR)"; \
		$(GO_LINT) $(TOP_PACKAGE_DIR)/$$p; \
	done

test:
	@for p in $(PACKAGE_LIST); do \
		echo "$(OK_COLOR)==> Unit Testing $$p ...$(NO_COLOR)"; \
		$(GO_TEST_VERBOSE) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

update-deps:
	@for p in $(PACKAGE_LIST); do \
		echo "$(OK_COLOR)==> Update dependencies for $$p ...$(NO_COLOR)"; \
		$(GO_DEPS_UPDATE) || exit 1; \
	done

vet:
	@for p in $(PACKAGE_LIST); do \
		echo "$(OK_COLOR)==> Vet $$p ...$(NO_COLOR)"; \
		$(GO_VET) $(TOP_PACKAGE_DIR)/$$p; \
	done


.PHONY: build clean deps fmt lint run test update-deps vet
