BIN_NAME=rss-godler

build-windows:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${BIN_NAME}.exe -ldflags "-s -w" . 

build-linux:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${BIN_NAME} -ldflags "-s -w" . 

build-darwin:
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${BIN_NAME}_darwin -ldflags "-s -w" . 

all: build-windows build-linux build-darwin