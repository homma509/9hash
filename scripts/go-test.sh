#!/usr/bin/env bash

sleep 1

PACKAGE=${1:-...}
ARGS=${2}
COVER=./test/cover

scripts/wait-dynamodb-local.sh

if [ $? -ne 0 ]; then
  echo "Failed to launch dynamo local"
  exit 1
fi

if [ ! -d $COVER ]; then
  mkdir -p $COVER
fi

pwd
echo "--------------"
echo "Running go test."
echo
richgo test ./$PACKAGE $ARGS -v -coverpkg=./$PACKAGE -coverprofile=$COVER/cover.out
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

if [ -r $COVER/cover.out ]; then
  go tool cover -html=$COVER/cover.out -o $COVER/cover.html
  echo "output coverage file to $COVER/cover.html."
fi
