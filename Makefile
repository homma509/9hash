DOCKER_YAML=-f docker-compose.yml
DOCKER=COMPOSE_PROJECT_NAME=9hash docker-compose $(DOCKER_YAML)

docker-build:
	$(DOCKER) build ${ARGS}

docker-up:
	$(DOCKER) up

go-lint:
	$(DOCKER) run go ./scripts/go-lint.sh

go-build:
	$(DOCKER) run go ./scripts/go-build.sh

npm-serve:
	$(DOCKER) run --service-ports npm npm run serve

npm-build:
	$(DOCKER) run npm npm run build

sample:
	$(DOCKER) run sls sls create --template aws-go --name serverless-sample --path ./serverless_sample

package:
	$(DOCKER) run sls sls package

deploy:
	$(DOCKER) run sls sls deploy

clean:
	./scripts/clean.sh

remove:
	$(DOCKER) run sls sls remove
