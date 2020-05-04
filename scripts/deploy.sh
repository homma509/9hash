#!/usr/bin/env bash

npm install

ENV="development"
STAGE="dev"
BACKEND=TRUE

while [ $# -gt 0 ]
do
    case "$1" in
        "dev" ) ENV="development"
                STAGE="dev"
                ;;
        "staging" ) ENV="staging"
                    STAGE="staging"
                    ;;
        "prod" ) ENV="production"
                 STAGE="prod"
                 ;;
        "frontend" ) BACKEND=FALSE
                     ;;
    esac
    shift
done

if [ $BACKEND = TRUE ]; then
    sls deploy --env $ENV --stage $STAGE --config serverless_go.yml
    API_URL=$(sls info --stage $STAGE --config serverless_go.yml --verbose | sed -n -r 's/^.*ServiceEndpoint: (.+)/\1/p')
    echo -n -e "NODE_ENV=$ENV\nVUE_APP_API_BASE_URL=$API_URL/v1" > ./web/app/.env
else
    sls deploy --env $ENV --stage $STAGE --config serverless_npm.yml
fi
