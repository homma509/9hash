#!/usr/bin/env bash

check() {
  echo "Wait for $1"
  try=0
  until curl -X POST --connect-timeout 10 --max-time 10 $1 > /dev/null 2>&1; do
    echo -n "."
    sleep 1
    try=$(expr $try + 1)
    if [ $try -ge 10 ]; then
      echo ""
      echo "Failed to wait for $1"
      exit 1
    fi
  done
}

check $DYNAMODB_LOCAL_ENDPOINT
