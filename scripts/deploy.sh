#!/usr/bin/env bash

npm install

SYSTEM=${1:-backend}
STAGE=${2:-dev}

case "$STAGE" in
    "dev" )     ENV="development"   ;;
    "staging" ) ENV="staging"       ;;
    "prod" )    ENV="production"    ;;
esac

if [ "$SYSTEM" = "backend" ]; then
    sls deploy --env $ENV --stage $STAGE --config serverless_go.yml
    API_URL=$(sls info --stage $STAGE --config serverless_go.yml --verbose | sed -n -r 's/^.*ServiceEndpoint: (.+)/\1/p')
    echo -n -e "NODE_ENV=$ENV\nVUE_APP_API_BASE_URL=$API_URL/v1" > ./web/app/.env
else
    sls deploy --env $ENV --stage $STAGE --config serverless_npm.yml
fi
