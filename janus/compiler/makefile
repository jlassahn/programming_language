
MAKEFILE_DIR := $(abspath $(dir $(firstword $(MAKEFILE_LIST))))

export GOPATH=$(MAKEFILE_DIR)/go_src

all:
	go build go_src/src/janusc.go

