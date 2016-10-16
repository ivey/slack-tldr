export COMMIT_SHA := $(shell git rev-parse --short HEAD)
ifdef TRAVIS_TAG
	VERSION := $(TRAVIS_TAG)
else
	VERSION := "git-$(COMMIT_SHA)"
endif

slack-tldr:
	CGO_ENABLED=0 go build -ldflags="-X main.Version=$(VERSION)"

build: clean test slack-tldr.linux.amd64.tgz

slack-tldr.linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o "$@" -ldflags="-X main.Version=$(VERSION)"

slack-tldr.linux.amd64.tgz: slack-tldr.linux
	mv slack-tldr.linux slack-tldr
	tar cvfz slack-tldr.linux.amd64.tgz slack-tldr slack-tldr.service
	rm slack-tldr

test:
	go test -v

clean:
	rm -f slack-tldr 
