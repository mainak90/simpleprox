UNAME := $(shell uname)

ifeq ($(UNAME), Linux)
	OARCH=linux
endif

ifeq ($(UNAME), Darwin)
	OARCH=darwin
endif

default: build

build:
	cp config.json bin/ && GOOS=${OARCH} GOARCH=amd64 go build -o bin/simpleprox server.go