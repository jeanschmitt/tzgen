.PHONY: build

build:
	go build -o ./build/tzgen -ldflags "-X 'main.Version=`git tag --sort=-version:refname | head -n 1`'" .
