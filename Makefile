TAG ?= latest

build-image:
	docker build -t hyha-api:$(TAG) .

demo:
	$(MAKE) build-image;
	@TAG=$(TAG) docker-compose -f dev/docker-compose.yaml ps | grep -q 'hypha-api' && \
        { TAG=$(TAG) docker-compose -f dev/docker-compose.yaml up -d --no-deps --build hypha-api; } || \
        { TAG=$(TAG) docker-compose -f dev/docker-compose.yaml up -d; }

demo-cleanup:
	docker-compose -f dev/docker-compose.yaml down