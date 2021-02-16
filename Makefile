.NOTPARALLEL: plugins

all: plugins

.PHONY: plugins
plugins:
	@echo "Building plugins..."
	@for pl in $(shell sh -c "ls */main.go"); do go build -ldflags="-s -w" --buildmode=plugin -o . $$PWD/$${pl::-7}; done;
