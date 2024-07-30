TAG ?= latest

build-image:
	docker build -t hyha-api:$(TAG) .

demo:
	$(MAKE) build-image;
	@TAG=$(TAG) docker-compose -f dev/docker-compose.yaml up -d --force-recreate;

demo-cleanup:
	docker-compose -f dev/docker-compose.yaml down