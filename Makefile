BUILD_TARGET ?=
APP_VERSION ?= 0.0.0-dev
DPKG_APP_VERSION := $(shell echo ${APP_VERSION} | sed 's/v//')-0
GIT_SHA ?= $(shell git log -1 --pretty=format:"%h")
APP_BUILD ?= $(shell date -u "+%Y%m%d-%H%M")-${GIT_SHA}
BUILD_BINARY ?= build/hk-svcs-bridge-${APP_VERSION}-${BUILD_TARGET}
UNAME_S ?= $(shell uname -s)

CI_API_V4_URL ?= https://gitlab.com/api/v4
CI_PROJECT_ID ?= 43209634
PKGS_REGISTRY ?= ${CI_API_V4_URL}/projects/${CI_PROJECT_ID}/packages/generic/${APP_VERSION}

CGO_ENABLED=0
GITLAB_AUTH_HEADER ?= "JOB-TOKEN: ${CI_JOB_TOKEN}"

ifeq($(CI_JOB_TOKEN),)
GITLAB_AUTH_HEADER ?= "PRIVATE-TOKEN: ${GITLAB_TOKEN}" 
endif

.EXPORT_ALL_VARIABLES:

DEB_ARCH ?= $(GOARCH)
ifeq ($(GOARCH),arm)
DEB_ARCH := armhf
endif

_DEB_BUILD_PATH := $(shell mktemp -d)/build/deb/${DEB_ARCH}/hk-svcs-bridge-${APP_VERSION}
DATE ?= $(shell date -R)

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
	$(MAKE) build -e BUILD_BINARY=/usr/bin/hk-svcs-bridge
	cp debian/hk-svcs-bridge.service /etc/systemd/system/hk-svcs-bridge.service
	cp config.yaml /etc/hk-svcs-bridge.yaml
	mkdir -p /usr/var/hk-svcs-bridge/
	systemctl daemon-reload
	systemctl enable /etc/systemd/system/hk-svcs-bridge.service
else
	@echo "Error: make install cmd supports only GNU/Linux"
endif
.PHONY: install

build:
	mkdir -p build
	go build \
		-ldflags "-X main.app_version=${APP_VERSION} -X main.app_build=${APP_BUILD}" \
		-o ${BUILD_BINARY}
	cp LICENSE build/
	cp config.yaml build/example-config.yaml
.PHONY: build

dpkg:
	mkdir -p ${_DEB_BUILD_PATH}/build
	cp -aR debian ${_DEB_BUILD_PATH}/
	cp LICENSE config.yaml ${_DEB_BUILD_PATH}/build
	cp ${BUILD_BINARY} ${_DEB_BUILD_PATH}/build/hk-svcs-bridge
	chmod +x ${_DEB_BUILD_PATH}/debian/rules
	envsubst < debian/control > ${_DEB_BUILD_PATH}/debian/control
	envsubst < debian/changelog > ${_DEB_BUILD_PATH}/debian/changelog
	cd ${_DEB_BUILD_PATH};\
	dpkg-buildpackage -us -uc -d --host-arch ${DEB_ARCH}
	cp ${_DEB_BUILD_PATH}/../*.deb ${_DEB_BUILD_PATH}/../*.dsc build/
.PHONY:

run: clean
	LOG_LEVEL=debug go run ./*.go
.PHONY: run

# Upload to GitLab Package Registry
upload_pkgs:
	for i in build/*; do \
		curl --header "${GITLAB_AUTH_HEADER}" --upload-file $$i "${PKGS_REGISTRY}/$$i"; \
	done
.PHONY: upload_pkgs