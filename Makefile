BINARY=kubeauth
VERSION=`cat VERSION`
BUILD=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)"

all:
	go build -o $(BINARY) $(LDFLAGS)

clean:
	rm $(BINARY)