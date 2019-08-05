# Note: git tags *must* be annotated for this to work
VERSION := $(shell git describe --tags `git rev-list --tags --max-count=1` | sed 's/^[^0-9]*//')

all: instafeed

instafeed:
	@go build -mod=vendor -ldflags "-s -w -X main.version=$(VERSION)" -o instafeed
 