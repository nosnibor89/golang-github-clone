
start-api: build start-local-api

invoke-function: build invoke

cli:
	cd ./src/github-clone-cli && go install

build:
	sam build

#invoke-debug:
#	sam local invoke -d 2345 --debugger-path . --debug-args="-delveAPI=2" -e $(event) --env-vars env.local.json

invoke:
	sam local invoke $(name) -e $(event) --env-vars env.local.json

start-local-api:
	sam local start-api --env-vars env.local.json
