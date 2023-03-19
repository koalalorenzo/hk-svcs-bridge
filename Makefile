BUILD_TARGET ?=
APP_VERSION ?= local-dev
GIT_TAG ?= $(shell git log -1 --pretty=format:"%h")
APP_BUILD ?= $(shell date -u "+%Y%m%d-%H%M")-${GIT_TAG}

ifeq (${BUILD_TARGET},rpi)
GOARCH := arm
GOOS := linux
GOARM=7
endif

CGO_ENABLED=0

.EXPORT_ALL_VARIABLES:

clean:
	rm -rf build
.PHONY: build

test: clean
	go fmt $(go list ./... | grep -v /vendor/)
	go vet $(go list ./... | grep -v /vendor/)
	go test -race $(go list ./... | grep -v /vendor/)	
.PHONY: test

build:
	mkdir -p build
	go build \
		-ldflags "-X main.app_version=${APP_VERSION} -X main.app_build=${APP_BUILD}" \
		-o build/go-hk-svc-${APP_VERSION}-${BUILD_TARGET}

run: clean
	LOG_LEVEL=debug go run ./*.go
.PHONY: run