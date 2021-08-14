.NOTPARALLEL: build
SHELL=sh

all: build
static: static-arg build

.PHONY: static-arg
static-arg:
	$(eval STATIC_LDFLAGS = -extldflags "-static")

.PHONY: build
build:
	@echo "Building plugins..."
	@for pl in $(shell sh -c "ls */main.go"); do go build -ldflags="-s -w $(STATIC_LDFLAGS)" --buildmode=plugin -o . $$PWD/$${pl%/main.go}; done;
