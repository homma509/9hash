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

npm-lint:
	$(DOCKER) run npm npm run lint

npm-build:
	$(DOCKER) run npm npm run build

# npm-clean:
# 	$(DOCKER) run npm npm cache clean --force

sample:
	$(DOCKER) run sls sls create --template aws-go --name serverless-sample --path ./serverless_sample

package:
	$(DOCKER) run sls sls package

deploy:
	$(DOCKER) run sls ./scripts/deploy.sh backend ${ARGS}
	$(DOCKER) run npm npm run build
	$(DOCKER) run sls ./scripts/deploy.sh frontend ${ARGS}
	# make ARGS="prod" deploy

clean:
	./scripts/clean.sh

remove:
	$(DOCKER) run sls ./scripts/remove.sh ${ARGS}
	# make ARGS="prod" remove
