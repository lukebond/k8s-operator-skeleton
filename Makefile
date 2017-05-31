REGISTRY := quay.io
NAMESPACE := lukebond
CONTAINER_NAME := k8s-operator-skeleton
VERSION := v1.0.0
CONTAINER_TAG = ${REGISTRY}/${NAMESPACE}/${CONTAINER_NAME}:$(VERSION)
IMAGE_FILENAME = $(CONTAINER_NAME).tar
GITREF  = $(shell git show --oneline -s | head -n 1 | awk '{print $$1}')

.PHONY: all build

all: build run

deps:
	vndr

build:
	CGO_ENABLED=0 go build -o bin/$(CONTAINER_NAME) cmd/operator/main.go
	docker build -t $(CONTAINER_TAG) --build-arg BUILDDATE=`date -u +%Y-%m-%dT%H:%M:%SZ` --build-arg VERSION=$(VERSION) --build-arg VCSREF=$(GITREF) .

save:
	docker save $(CONTAINER_TAG) -o $(IMAGE_FILENAME)
	chmod +r $(IMAGE_FILENAME)

clean:
	docker rmi $(CONTAINER_TAG)

clean-minikube:
	kubectl delete deployment k8s-operator-skeleton 2>/dev/null || true
	kubectl delete thirdpartyresource example.lukeb0nd.com 2>/dev/null || true
