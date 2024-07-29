CONTAINER_CMD=$(shell which docker || which podman 2>/dev/null)
TAG ?= latest
BUILD_DEMO ?= true

build-image:
	$(CONTAINER_CMD) build -t hyha-api:$(TAG) .

demo:
	[ "$(BUILD_DEMO)" = "true" ] && $(MAKE) build-image
	$(CONTAINER_CMD) network create hypha-network
	$(CONTAINER_CMD) run -d --name postgres-hypha --network hypha-network -e POSTGRES_DB=hypha -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=mysecretpassword postgres
	sleep 10
	$(CONTAINER_CMD) run -d --name hypha-api --network hypha-network -p 8081:8081 -v $(PWD)/dev/config.yaml:/config.yaml hyha-api:$(TAG) hypha-api --config /config.yaml

demo-cleanup:
	$(CONTAINER_CMD) stop hypha-api
	$(CONTAINER_CMD) stop postgres-hypha
	$(CONTAINER_CMD) rm hypha-api
	$(CONTAINER_CMD) rm postgres-hypha
	$(CONTAINER_CMD) network rm hypha-network