
start-api: build start-local-api

cli:
	cd ./src/github-clone-cli && go install

build:
	sam build

start-local-api:
	sam local start-api --env-vars env.local.json
