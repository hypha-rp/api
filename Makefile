TAG ?= latest

build-image:
	docker build --no-cache -t hypha-api:$(TAG) .

#######################################
### Development Environment Targets ###
#######################################

dev-up: build-image
	@TAG=$(TAG) docker-compose -f dev-tools/docker/api-compose.yaml up -d --force-recreate;

dev-down:
	docker-compose -f dev-tools/docker/api-compose.yaml down

dev-test: dev-down dev-up dev-product-create-and-report

dev-product-create-and-report:
	sleep 10
	./dev-tools/scripts/create-product-report-results.sh