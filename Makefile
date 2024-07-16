lint:
	cd cart && make lint

coverage:
	cd cart && make coverage

build-all:
	cd cart && GOOS=linux GOARCH=amd64 make build

run-all: build-all
	docker-compose up --force-recreate --build
