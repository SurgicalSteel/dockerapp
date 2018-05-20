#!/bin/bash

build:
	@env GOOS=linux GOARCH=amd64 go build

run:
	@dockerapp