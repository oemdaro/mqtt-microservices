#
# Bootstrap command to compile protobuf
#
SHELL := /bin/bash

compile:
	@protoc --go_out=plugins=grpc:. *.proto

install:
	@go install
