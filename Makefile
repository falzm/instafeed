all: instafeed

instafeed:
	@GOPATH=$(PWD)/vendor go build ./src/cmd/instafeed
