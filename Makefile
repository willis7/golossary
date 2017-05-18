NAME = golossary
PWD := $(MKPATH:%/Makefile=%)

help:
	@echo "Usage:"
	@echo "    make <target>"
	@echo
	@echo "Available targets: "
	@echo "    build                - performs a full build of the project (clean install check)"
	@echo "    compile				- creates a binary in bin directory of GOPATH"
	@echo "    check                - performs all verification tasks in the project"
	@echo "    coverage             - print a coverage report to terminal"
	@echo "    clean                - deletes the project vendor directory."
	@echo "    install              - download all dependencies"
	@echo "    lint                 - ensure code is standards compliant"
	@echo "    test            		- run tests"
	@echo "    docker-build         - build docker image"
	@echo "    docker-run           - run docker conatiner"
	@echo "    docker-rm            - remove docker container"
	@echo "    docker-rmi           - remove docker image"
	@echo


build:	clean install compile check

check:	test

clean :
	cd "$(PWD)"
	rm -rf vendor

compile:
	go install main.go

coverage:
	echo 'mode: atomic' > coverage.txt && go list $(shell glide novendor) | xargs -n1 -I{} sh -c 'go test -covermode=atomic -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> coverage.txt' && rm coverage.tmp

fmt:
	go fmt ./...

test:
	go test -v $(shell glide novendor)

race:
	go test -race -v $(shell glide novendor)

run:
	go run main.go

install:
	glide install

docker-build:
	GOOS=linux go build -o app
	docker build -t willis7/$(NAME) .
	rm -f app

docker-run:
	docker run -it --rm --name $(NAME) willis7/$(NAME)

docker-rm:
	docker rm $(NAME)

docker-rmi:
	docker rmi $(NAME)

default: help

