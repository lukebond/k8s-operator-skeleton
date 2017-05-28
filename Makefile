REGISTRY := quay.io
NAMESPACE := lukebond
CONTAINER_NAME := k8s-operator-skeleton
VERSION := v1.0.0
CONTAINER_TAG = ${REGISTRY}/${NAMESPACE}/${CONTAINER_NAME}:$(VERSION)
GITREF  = $(shell git show --oneline -s | head -n 1 | awk '{print $$1}')

.PHONY: all build

all: build run

deps:
	vndr

build:
	CGO_ENABLED=0 go build -o bin/$(CONTAINER_NAME) cmd/operator/main.go
	docker build -t $(CONTAINER_TAG) --build-arg BUILDDATE=`date -u +%Y-%m-%dT%H:%M:%SZ` --build-arg VERSION=$(VERSION) --build-arg VCSREF=$(GITREF) .

clean:
	docker rmi $(CONTAINER_TAG)
