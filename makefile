test: install-static-check
	go test -failfast -count=1 -race -v ./...
	staticcheck -checks=all ./...

sleep-one-sec:
	sleep 1;

install-static-check:
	go install honnef.co/go/tools/cmd/staticcheck@latest

db-local:
	docker compose -f dev.docker-compose.yaml up -d database

migrate-local:
	docker compose -f dev.docker-compose.yaml up migrate

migrate-down:
	docker compose -f dev.docker-compose.yaml run --rm migrate down ${step}

setup: db-local sleep-one-sec migrate-local

change-db:
	docker compose -f dev.docker-compose.yaml run --no-deps --rm migrate create --dir=migrations --ext=sql --seq ${name}

run:
	go run app/api/main.go

tear:
	docker compose -f dev.docker-compose.yaml down

run-local: setup
	go run app/api/main.go

test-local: setup test

e2e:
	go run app/salvo/main.go

slow-e2e:
	sleep 7
	make e2e

e2e-local:
	make slow-e2e &
	make run-local

tidy:
	go mod tidy
	go mod vendor

update-dependencies:
	go get -u ./...
	go mod tidy
	go mod vendor

coverage:
	go test -coverprofile=cp.out ./...
	go tool cover -html=cp.out

coverage-local: setup coverage
