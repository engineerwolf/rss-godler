BIN_NAME=rss-godler
VERSION=$(shell git describe --always --long --dirty)
BUILD_NUMBER=$(shell date +%d%m%Y_%H%M%S)

build-windows:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${BIN_NAME}.exe -ldflags "-s -w -X main.version=${VERSION}.${BUILD_NUMBER}" . 

build-linux:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${BIN_NAME} -ldflags "-s -w -X main.version=${VERSION}.${BUILD_NUMBER}" . 

build-darwin:
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${BIN_NAME}_darwin -ldflags "-s -w -X main.version=${VERSION}.${BUILD_NUMBER}" . 

all: build-windows build-linux build-darwin