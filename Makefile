.ONESHELL:
SHELL = /bin/bash

BINARY_NAME = semtag
CHANGELOG_NAME = CHANGELOG.md



.PHONY: help
help: ## show this help
	@echo -e '\nUsage: make [target] ...\n'
	# echo -e "TARGET: DEPENDENCIES ## DESCRIPTION \n`egrep '^[0-9a-zA-Z](.*)+:.*?## ' $(MAKEFILE_LIST)`" | awk 'NR ==1 {print $0} ; NR > 1 {print $0 | "sort"}' | column -t -c 2 -s ':#'
	echo -e "TARGET: DEPENDENCIES ## DESCRIPTION \n`egrep '^[0-9a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST)`" | awk 'NR ==1 {print $0} ; NR > 1 {print $0 | "sort"}' | column -t -c 2 -s ':#'

dependencies: ## download the module dependencies
	@set -euo pipefail
	@echo -e "\n\tdownload the module dependencies"
	go mod download -x

compile: dependencies ## compile the packages and the dependencies
	@set -euo pipefail
	@echo -e "\n\tcode format for packages and dependencies"
	go fmt ./...
	@echo -e "\n\tcompile packages and dependencies"
	go build ./...

test: dependencies ## run tests
	@set -euo pipefail
	@echo -e "\n\texecute tests"
	go test ./... -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html



build: build-freebsd build-macos build-linux build-windows ## create binaries

build-freebsd: compile ## create binaries for FreeBSD
	GOOS=freebsd GOARCH=386 go build -o bin/$(BINARY_NAME)-freebsd-386 main.go
	GOOS=freebsd GOARCH=amd64 go build -o bin/$(BINARY_NAME)-freebsd-amd64 main.go

build-macos: compile ## create binaries for MacOS
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 main.go

build-linux: compile ## create binaries for Linux
	GOOS=linux GOARCH=386 go build -o bin/$(BINARY_NAME)-linux-386 main.go
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 main.go
	cp -v bin/$(BINARY_NAME)-linux-amd64 bin/$(BINARY_NAME)

build-windows: compile ## create binaries for Windows
	GOOS=windows GOARCH=386 go build -o bin/$(BINARY_NAME)-windows-386 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64 main.go



docker: ## create a Docker image
	@set -euo pipefail
	@echo -e "\n\tcreate the Docker image"
	docker build --tag semtag .

tag: build ## create an annotated git tag
	@set -euo pipefail
	git config --list | grep "user.email" || git config --global user.email "ci@foo.com" && git config --global user.name "ci"
	bin/semtag -increment=auto -git-tag -prefix="v" -push -path=main.go -path=go.mod -path=internal -path=pkg

.PHONY: changelog
changelog: ## create the repository changelog
	@set -euo pipefail
	@echo -e "\n\tcreate the change log"
	GIT_COMMIT_URL="https://github.com/mpdred/semantic-tagger/commit/" \
	GIT_TAG_URL="https://github.com/mpdred/semantic-tagger/releases/tag/" \
	bin/semtag -changelog -prefix=v


.PHONY: clean
clean: # deletes all temporary files
	rm -rfv \
		c.out coverage.html \
		.version-* \
		bin/ \
		$(CHANGELOG_NAME) \
	|| true
