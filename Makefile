BIN_NAME=pubsub
VERSION=1.0.0

all:
	$(MAKE) build-linux build-darwin
	$(MAKE) upx-linux upx-darwin
	$(MAKE) sha256-linux sha256-darwin

mkdir-base:
	mkdir -p "bin/$(VERSION)/$(GOOS)/$(GOARCH)"

mkdir-linux:
	$(MAKE) mkdir-base GOOS=linux GOARCH=amd64

mkdir-darwin:
	$(MAKE) mkdir-base GOOS=darwin GOARCH=amd64

build-base:
	go build -ldflags="-w -s -X grasys/pubsub/cmd.version=$(VERSION)" -trimpath -o bin/$(VERSION)/$(GOOS)/$(GOARCH)/$(BIN_NAME) main.go

build-linux:
	$(MAKE) build-base GOOS=linux GOARCH=amd64

build-darwin:
	$(MAKE) build-base GOOS=darwin GOARCH=amd64

upx-install:
	@if ! type upx; then brew install upx ; fi

upx-base:
	$(MAKE) upx-install
	upx --best --lzma bin/$(VERSION)/$(GOOS)/$(GOARCH)/$(BIN_NAME)

upx-linux:
	$(MAKE) upx-base GOOS=linux GOARCH=amd64

upx-darwin:
	$(MAKE) upx-base GOOS=darwin GOARCH=amd64

sha256-base:
	shasum -a 256 bin/$(VERSION)/$(GOOS)/$(GOARCH)/$(BIN_NAME)

sha256-linux:
	$(MAKE) sha256-base GOOS=linux GOARCH=amd64

sha256-darwin:
	$(MAKE) sha256-base GOOS=darwin GOARCH=amd64
