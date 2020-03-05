.PHONY: init build deploy

init:
	@type sls > /dev/null || npm install -g serverless@1.63.0
	go get -v ./...

build:
	GOOS=linux go build -o bin/serverless ./

deploy: build
	sls deploy -v
