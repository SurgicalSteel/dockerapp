#!/bin/bash

build:
	@go build

build-docker:
	@env GOOS=linux GOARCH=amd64 go build

run:
	@dockerapp

build-image:
	@docker build --tag ceruntu .