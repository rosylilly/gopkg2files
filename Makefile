GO?=go

GO_ALL_PACKAGES:=$(shell ${GO} list ./...)
GO_ROOT_PACKAGE:=$(shell ${GO} list .)
GO_CMD_PACKAGES:=$(shell ${GO} list ./cmd/...)

GO_CMD_ARTIFACTS:=$(patsubst ${GO_ROOT_PACKAGE}/cmd/%,bin/%,${GO_CMD_PACKAGES})

GOPKG2FILES?=${GO} run ./cmd/gopkg2files

.PHONY: build
build: ${GO_CMD_ARTIFACTS}

.SECONDEXPANSION:
bin/%: Makefile $$(shell ${GOPKG2FILES} ${GO_ROOT_PACKAGE}/cmd/%)
	@${GO} build -o $@ $(patsubst bin/%,${GO_ROOT_PACKAGE}/cmd/%,$@)
