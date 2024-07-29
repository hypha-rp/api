CONTAINER_CMD=$(shell which docker || which podman 2>/dev/null)
TAG ?= latest

build-container:
	$(CONTAINER_CMD) build -t hyha-api:$(TAG) .
