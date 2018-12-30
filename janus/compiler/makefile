
MAKEFILE_DIR := $(abspath $(dir $(firstword $(MAKEFILE_LIST))))

export GOPATH=$(MAKEFILE_DIR)/go_src

all:
	go build go_src/src/janusc.go
	gcc -c library/clib/clib.c -o library/clib/clib.o

test:
	go test -v -coverpkg driver,output,parser,symbols,generator  -coverprofile=coverage.out driver
	go tool cover -html=coverage.out -o coverage.html

clean:
	-rm tests/*.out
	-rm tests/*.ll
	-rm tests/*.s
	-rm library/clib/clib.o

