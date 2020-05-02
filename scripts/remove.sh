#!/usr/bin/env bash

npm install

if [ $# -eq 0 ]; then
    ENV="development"
    STAGE="dev"
elif [ $# -eq 1 ]; then
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
    esac
fi

sls remove --env $ENV --stage $STAGE