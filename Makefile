.DEFAULT_GOAL := build

BIN_DIR=bin
BIN_FILE=SearchAndSortMovies
VERSION_FILE=VERSION
GO_BIN=$(shell which go)

HAS_GO_BIN := $(shell command -v go 2> /dev/null)

GIT_STATUS=$(shell git status --porcelain)
BUILD_VERSION:=$(shell git log --pretty=format:'%h' -n 1)
BUILD_DATE:=$(shell date '+%Y-%m-%d_%k:%M:%S')

ifneq ($(wildcard $(VERSION_FILE)),)
	VERSION:=$(shell cat $(VERSION_FILE))
else
	VERSION:=
endif

ifeq ($(GIT_STATUS),)
  BUILD_CLEAN=yes
else
  BUILD_CLEAN=no
endif

.PHONY: check-env
check-env:
ifndef GOPATH
	$(error GOPATH is undefined, you need to configure your system ...)
endif
# ifndef GOROOT
# 	$(error GOROOT is undefined, you need to configure your system ...)
# endif
ifndef HAS_GO_BIN
	$(error go command couldn't be found in PATH)
endif

bin:
	mkdir ${BIN_DIR}

.PHONY: build
build: ${BIN_DIR} check-env
	${GO_BIN} build -ldflags "-X 'main.BuildVersion=${VERSION}' -X 'main.BuildHash=${BUILD_VERSION}' -X 'main.BuildDate=${BUILD_DATE}' -X 'main.BuildClean=${BUILD_CLEAN}'" -o "${BIN_DIR}/${BIN_FILE}" .

.PHONY: clean
clean:
	rm -fv ${BIN_DIR}/${BIN_FILE}