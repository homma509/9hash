#!/usr/bin/env bash

sleep 1

PACKAGE=${1:-...}
ARGS=${2}

scripts/wait-dynamodb-local.sh

if [ $? -ne 0 ]; then
  echo "Failed to launch dynamo local"
  exit 1
fi

pwd
echo "--------------"
echo "Running go test."
echo
richgo test ./$PACKAGE $ARGS -v
RETURNCD=$?
if [ $RETURNCD -ne 0 ]; then
  echo "--------------"
  echo "go test FAILED."
  echo "--------------"
  echo "RETURNCD: $RETURNCD"
  exit $RETURNCD
else
  echo "go test OK."
fi

echo "--------------"
echo "TEST DONE."
echo "--------------"
