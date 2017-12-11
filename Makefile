SHELL := /bin/bash

.PHONY: all

all: image

image: auth-service-image data-service-image mqtt-image

auth-service-image:
	$(MAKE) -C auth-service

data-service-image:
	$(MAKE) -C data-service

mqtt-image:
	docker build -f mqtt-server/Dockerfile -t local/mqtt-server .
