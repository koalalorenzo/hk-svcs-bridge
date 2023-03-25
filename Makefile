BUILD_TARGET ?=
APP_VERSION ?= local-dev
GIT_SHA ?= $(shell git log -1 --pretty=format:"%h")
APP_BUILD ?= $(shell date -u "+%Y%m%d-%H%M")-${GIT_SHA}
BUILD_BINARY ?= build/hk-systemd-bridge-${APP_VERSION}-${BUILD_TARGET}
UNAME_S ?= $(shell uname -s)

CGO_ENABLED=0

ifeq ($(GOARCH),arm)
	DEB_ARCH := armhf
endif
DEB_ARCH ?= $(GOARCH)

.EXPORT_ALL_VARIABLES:

clean:
	rm -rf build
.PHONY: build

test: clean
	go fmt $(go list ./... | grep -v /vendor/)
	go vet $(go list ./... | grep -v /vendor/)
	CGO_ENABLED=1 go test -race $(go list ./... | grep -v /vendor/)	
.PHONY: test

install: clean
ifeq ($(UNAME_S),Linux)
	$(MAKE) build -e BUILD_BINARY=/usr/bin/go-hk-systemd-bridge
	cp systemd.service /etc/systemd/system/go-homekit-systemd-bridge.service
	cp config.yaml /etc/go-homekit-systed-bridge.yaml
	mkdir -p /usr/var/go-hk-systemd-bridge/
	systemctl daemon-reload
	systemctl enable /etc/systemd/system/go-homekit-systemd-bridge.service
else
	@echo "Error: make install cmd supports only GNU/Linux"
endif
.PHONY: install

build:
	mkdir -p build
	go build \
		-ldflags "-X main.app_version=${APP_VERSION} -X main.app_build=${APP_BUILD}" \
		-o ${BUILD_BINARY}

run: clean
	LOG_LEVEL=debug go run ./*.go
.PHONY: run