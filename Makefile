.ONESHELL:
SHELL = /bin/bash

all: | test build

build:
	set -e
	go build -o /tmp/semtag

test:
	set -e
	go test ./... -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html

tag:
	set -e
	git config --list | grep "user.email" || git config --global user.email "ci@foo.com" && git config --global user.name "ci"
	.//tmp/semtag -increment=auto -git-tag -prefix="v" -push -path=main.go -path=go.mod -path=internal -path=pkg

changelog:
	set -e
	GIT_COMMIT_URL="https://github.com/mpdred/semantic-tagger/commit/" \
	GIT_TAG_URL="https://github.com/mpdred/semantic-tagger/releases/tag/" \
	/tmp/semtag -changelog -prefix=v -changelog-regex="^v[0-9]+\.[0-9]+\.[0-9]+$$"


upload:
	set -e
	cp -v /tmp/semtag /tmp/semtag-$$(git describe --tags `git rev-list --tags --max-count=1` | cut -d '.' -f1)
	sudo apt-get install -y awscli > /dev/null
	aws s3 sync /tmp/ s3://mpdred-public
