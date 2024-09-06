#CGO_ENABLED=0 go run -ldflags="-s -w -X 'main.production=${PRODUCTION}'" cmd/*.gorun:

.PHONY: run_auth
run_auth:
	CGO_ENABLED=0 go run auth/cmd/*.go -prod=false | go run common/tooling/logfmt/main.go