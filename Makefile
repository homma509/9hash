DOCKER_YAML=-f docker-compose.yml
DOCKER=COMPOSE_PROJECT_NAME=hash docker-compose $(DOCKER_YAML)

docker-build:
	$(DOCKER) build ${ARGS}

docker-up:
	$(DOCKER) up

go-lint:
	$(DOCKER) run go ./scripts/go-lint.sh

go-build:
	$(DOCKER) run go ./scripts/go-build.sh

sample:
	$(DOCKER) run sls sls create --template aws-go --name serverless-sample --path ./serverless_sample

package:
	$(DOCKER) run sls sls package

deploy:
	$(DOCKER) run sls ./scripts/deploy.sh

remove:
	$(DOCKER) run sls sls remove
