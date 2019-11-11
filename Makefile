PKGS := $(shell go  list ./... | grep -v /vendor)

uruk:
	CGO_ENABLED=0 go build -o bin/uruk ./pkg/main/

.PHONY: uruk_stripped
uruk_stripped:
	go build -o bin/uruk -ldflags="-s -w" ./pkg/main/

.PHONY: uruk_compressed
uruk_compressed: uruk_stripped
	upx bin/uruk