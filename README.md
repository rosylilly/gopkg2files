# gopkg2files

Convert go packages to files list.

## Usage

Install `github.com/rosylilly/gopkg2files/cmd/gopkg2files`.

Write your makefile like below:

```Makefile
GO?=go

GO_ALL_PACKAGES:=$(shell ${GO} list ./...)
GO_ROOT_PACKAGE:=$(shell ${GO} list .)
GO_CMD_PACKAGES:=$(shell ${GO} list ./cmd/...)

GO_CMD_ARTIFACTS:=$(patsubst ${GO_ROOT_PACKAGE}/cmd/%,bin/%,${GO_CMD_PACKAGES})

GOPKG2FILES?=gopkg2files

.PHONY: build
build: ${GO_CMD_ARTIFACTS}

.SECONDEXPANSION:
bin/%: Makefile $$(shell ${GOPKG2FILES} ${GO_ROOT_PACKAGE}/cmd/%)
	@${GO} build -o $@ $(patsubst bin/%,${GO_ROOT_PACKAGE}/cmd/%,$@)
```

and Run `make build` to build binaries on changed dependent go files.

## Options

```
$ ./bin/gopkg2files -h
Usage of gopkg2files:
  -cgo
    	list cgo files
  -debug
    	debug mode
  -goroot
    	list undered $GOROOT files
  -test
    	list test files
  -w string
    	Working directory (default ".")
  -xtest
    	list xtest files
```

## Author

Sho Kusano / @rosylilly
