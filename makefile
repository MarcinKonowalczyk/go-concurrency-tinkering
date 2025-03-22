

# Dynamically create a target for each directory containing tinkering.go
MAIN = tinkering.go
DIRS = $(wildcard ./*/$(MAIN))
DIRS := $(subst ./,,$(DIRS))
DIRS := $(subst /$(MAIN),,$(DIRS))
# $(info DIRS: $(DIRS))

all: $(DIRS)

define make-directory-target
.PHONY: $(1)
$(1):
	@echo "========== \033[1;32mRunning $(1) \033[0m =========="
	go run -C ./$(1) $(MAIN)

$(1)-race:
	@echo "========== \033[1;32mRunning $(1) \033[0m =========="
	go run -race -C ./$(1) $(MAIN)
endef

$(foreach element,$(DIRS),$(eval $(call make-directory-target,$(element))))

fmt:
	gofmt -w .