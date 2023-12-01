GO=$(shell which go)
LDFLAGS="-X \"main.versionInfo=$(shell cat version.txt)-${PIPELINE_IID}-$(shell date +"%Y.%m.%d")-$(shell go version | cut -d ' ' -f3)\""

.PHONY: test

build:
	$(GO) build -ldflags $(LDFLAGS) -mod=vendor -o $(GOPATH)/bin/yml2cdb cmd/yml2cdb/*
	$(GO) build -ldflags $(LDFLAGS) -mod=vendor -o $(GOPATH)/bin/yml2onlineconf cmd/yml2onlineconf/*

test:
ifdef t
	go test -v -count=1 -run $(t) ./...
else
	go test -v -count=1 ./...
endif
