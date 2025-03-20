

# Dynamically create a target for each directory containing tinkering.go

DIRS = $(wildcard ./*/tinkering.go)
DIRS := $(subst ./,,$(DIRS))
DIRS := $(subst /tinkering.go,,$(DIRS))
# $(info DIRS: $(DIRS))

all: $(DIRS)

define make-directory-target
.PHONY: $(1)
$(1):
	@echo "========== \033[1;32mRunning $(1) \033[0m =========="
	go run ./$(1)/tinkering.go

$(1)-race:
	@echo "========== \033[1;32mRunning $(1) \033[0m =========="
	go run -race ./$(1)/tinkering.go
endef

$(foreach element,$(DIRS),$(eval $(call make-directory-target,$(element))))

fmt:
	gofmt -w .