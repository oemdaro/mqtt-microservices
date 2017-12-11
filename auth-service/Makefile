SHELL := /bin/bash

godep := $(shell command -v dep 2> /dev/null)

.PHONY: all

all: dep image

dep:
ifndef godep
	$(error "dep is not available please install github.com/golang/dep")
endif
	dep ensure

image:
	docker build -t local/auth-service .
