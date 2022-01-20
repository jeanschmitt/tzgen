.PHONY: build

build:
	go build -o ./build/tzgen -s -ldflags "-X main.Version=`git tag --sort=-version:refname`" .
