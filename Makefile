DOCKER_YAML=-f docker-compose.yml
DOCKER=COMPOSE_PROJECT_NAME=hash docker-compose $(DOCKER_YAML)

build:
	$(DOCKER) build ${ARGS}

go-lint:
	$(DOCKER) run go-test ./scripts/go-lint.sh