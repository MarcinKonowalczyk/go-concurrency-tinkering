


all:
	@echo "No target specified. Please run 'make <target>' where <target> is one of the directories containing tinkering.go."

# Dynamically create a target for each directory containing tinkering.go

DIRS = $(wildcard ./*/tinkering.go)
DIRS := $(subst ./,,$(DIRS))
DIRS := $(subst /tinkering.go,,$(DIRS))
# $(info DIRS: $(DIRS))

define make-directory-target
.PHONY: $(1)
$(1):
	go run ./$(1)/tinkering.go
endef

$(foreach element,$(DIRS),$(eval $(call make-directory-target,$(element))))
