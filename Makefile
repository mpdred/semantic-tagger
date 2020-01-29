.ONESHELL:
SHELL = /bin/bash

all: | test build

build:
	go build -o out/semtag

test:
	go test ./... -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html

tag:
	git config --list | grep "user.email" || git config --global user.email "ci@foo.com" && git config --global user.name "ci"
	./out/semtag -git-tag -prefix v

upload:
	cp -v out/semtag out/semtag-$$(git describe --tags `git rev-list --tags --max-count=1` | cut -d '.' -f1)
	sudo apt-get install -y awscli > /dev/null
	aws s3 sync out/ s3://mpdred-public
