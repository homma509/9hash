# 9hash

9hash is a service to shorten URL using AWS Lambda and DynamoDB.

## Requirement

- Go 1.14

## Install

```
$ git clone https://github.com/homma509/9hash
```

## Test

backend all test

```
$ make go-test
```

backend package test 

```
$ make go-test PACKAGE=adapter/controller 
```

backend func test 

```
$ make go-test PACKAGE=adapter/controller ARGS="--run TestGetHash"
```

## build

backend
```
$ make go-build
```

frontend
```
$ make npm-build
```

## Deploy

development

```
$ make deploy
or
$ make deploy ARGS=dev
```

staging
```
$ make deploy ARGS=staging
```

production
```
$ make deploy ARGS=prod
```