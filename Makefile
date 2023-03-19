BUILD_TARGET ?=
APP_VERSION ?= local-dev
APP_BUILD ?= $(shell date -u "+%Y%m%d-%H%M")

ifeq (${BUILD_TARGET},rpi)
GOARCH := arm
GOOS := linux
GOARM=7
endif

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
		-o build/go-hk-svc-${GOARCH}-${GOOS}

run: clean
	LOG_LEVEL=debug go run ./*.go
.PHONY: run