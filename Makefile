PKG := github.com/jiradeto/gh-scanner
PKG_LIST := $(shell go list ${PKG}/...)
GOLINT?=		go run golang.org/x/lint/golint

setup:
	docker network create gh-scanner
start-service:
	docker-compose -f docker-compose-service.yaml up -d
stop-service:
	docker-compose -f docker-compose-service.yaml down
start-app-build:
	docker-compose -f docker-compose-app.yaml up -d --build
start-service-build:
	docker-compose -f docker-compose-service.yaml up -d --build
start-app:
	docker-compose -f docker-compose-app.yaml up -d
stop-app:
	docker-compose -f docker-compose-app.yaml down
start: start-service start-app-build

lint: 
	@echo $(GOLINT) -set_exit_status ${PKG_LIST}

test:
	@go test ${PKG_LIST}

mock/all:
	make mock/usecases m=repository
	make mock/repos m=repository

mock/usecases:
	mockgen \
		-source=./app/usecases/$(m)/main.go \
		-destination=./app/usecases/$(m)/mocks/$(m).go \
		-package $(m)usecasemocks \
        -mock_names UseCase=Mocks

mock/repos:
	mockgen \
		-source=./app/infrastructure/repos/$(m)/main.go \
		-destination=./app/infrastructure/repos/$(m)/mocks/$(m).go \
		-package $(m)repomocks \
        -mock_names Repo=Mocks
