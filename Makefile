.ONESHELL:
SHELL = /bin/bash

build:
	go build -o out/semtag

test:
	go test ./... -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html

tag:
	./out/semtag -tag git -prefix v

upload:
	cp -v out/semtag out/semtag-$$(git describe --tags `git rev-list --tags --max-count=1` | cut -d '.' -f1)
	sudo apt-get install -y awscli > /dev/null
	aws s3 sync out/ s3://mpdred-public
