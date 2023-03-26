CGO_ENABLED=1
PKG=github.com/fawad-khalil/go-user-service
GO=go
GOFLAGS=-v

default: build

build:
	go build

.PHONY: default build
