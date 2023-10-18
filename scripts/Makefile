.DEFAULT_GOAL := build

shadow:
	shadow ./../... # this tool detects shadowing variables
.PHONY:shadow

fmt:
	go fmt ./../...
.PHONY:fmt

lint: fmt
	golint ./../...
.PHONY:lint

vet: fmt
	go vet ./../...
.PHONY:vet

build: vet
	go build ../cmd/froxy
.PHONY:build