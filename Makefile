#CGO_ENABLED=0 go run -ldflags="-s -w -X 'main.production=${PRODUCTION}'" cmd/*.gorun:

.PHONY:build_auth
build_auth:
	CGO_ENABLED=0 go build -o bin/auth ./auth/cmd/main.go

.PHONY: run_auth
run_auth: build_auth
	bin/auth -prod=false | go run common/tooling/logfmt/main.go

.PHONY: compose_up
compose_up:
	docker-compose -f ./zarf/compose/docker-compose.dev.yml up --build --abort-on-container-exit

.PHONY: compose_down
compose_down:
	docker-compose -f ./zarf/compose/docker-compose.dev.yml down