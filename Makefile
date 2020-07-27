.PHONY: build tools dep bin clean aiappd package part publish test bench vet lint fmt vendor

export GO111MODULE=on

SHELL := /bin/bash
timestamp := $(shell date +"%s")
version := $(strip $(timestamp))
build_time := $(shell date +"%Y%m%d.%H%M%S")
build_sha := $(shell git rev-parse --verify HEAD)
goos = $(shell go env GOOS)

repo_name = $(shell basename `pwd`)
pkg_name = github.com/Arnold1/$(repo_name)
ldflags = -ldflags "-X $(pkg_name)/build.time=$(build_time) -X $(pkg_name)/build.number=$(version) -X $(pkg_name)/build.sha=$(build_sha)"
build_command = GOOS=$(goos) go build $(ldflags) -o
list = go list ./... | grep -v vendor/ | grep -v tools/

short_sha := $(shell git rev-parse --short=7 $(shell git rev-parse --verify HEAD))
binary_name ?= aiappd
local_img ?= $(binary_name):local
latest_sha_img ?= $(binary_name):$(short_sha)
latest_img ?= $(binary_name):latest
artifactory_tag ?= docker.demoorg.com/arnold1/ai-app/$(latest_sha_img)
artifactory_latest_tag ?= docker.demoorg.com/arnold1/ai-app/$(latest_img)
docker_location ?= ./provisioning/$(binary_name)/docker/
deploy_location ?= ./provisioning/$(binary_name)/deploy/
view_name ?= view

build: bin aiappd

bin:
	mkdir -p bin

clean:
	go clean -mod vendor ./...
	rm -rf ./bin

aiappd: bin
	$(build_command) bin/$(binary_name) $(pkg_name)/$(binary_name)
	cp -f ./bin/$(binary_name) $(docker_location)/$(binary_name)
	cp -r ./aiappd/$(view_name) $(docker_location)/$(view_name)
	pushd $(docker_location) && docker build -t $(local_img) . && rm -rf $(binary_name) && rm -rf $(view_name) && popd

locald: bin
	$(build_command) bin/locald $(pkg_name)/locald

run: aiappd
	docker run -p 8080:8080 --name=$(binary_name) --rm -it $(local_img)

publish:
	docker image tag $(local_img) $(artifactory_tag)
	docker image tag $(local_img) $(artifactory_latest_tag)
	#docker login -u "$(ARTIFACTORY_USERNAME)" -p "$(ARTIFACTORY_PASSWORD)" docker.blabla.com
	#docker push $(artifactory_tag)
	#docker push $(artifactory_latest_tag)
	#docker logout docker.blabla.com

publish-notag:
	docker image tag $(local_img) $(artifactory_tag)
	#docker login -u "$(ARTIFACTORY_USERNAME)" -p "$(ARTIFACTORY_PASSWORD)" docker.blabla.com
	#docker push $(artifactory_tag)
	#docker logout docker.blabla.com

test: fmt vet
	go test -race -mod vendor -p 1 -failfast -v ./... -coverprofile cover.out && \
	go tool cover -html=cover.out -o cover.html

bench:
	go test -mod vendor -p 1 -bench=. -run=xxx -benchmem=true -count=8 ./...

vendor:
	go mod vendor

vet:
	go vet -mod vendor ./...

fmt:
	go fmt ./...

golint := $(shell command -v golint)

lint: fmt
ifndef golint
	go get -u github.com/golang/lint/golint
endif
	$(list) | xargs golint | grep -v 'should have comment' || true