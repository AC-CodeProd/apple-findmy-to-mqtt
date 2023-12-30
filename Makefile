ENV_FILE_DEFAULT=.env

ifeq ($(strip $(ENV_FILE)),)
  #  $(error Path to .env file must be provided. Use -e option.)
else
  ENV_FILE_DEFAULT = $(ENV_FILE)
endif

-include $(ENV_FILE_DEFAULT)
export
.DEFAULT_GOAL := help

NAME=apple-findmy-to-mqtt
ROOT_DIR=$(shell pwd)
BUILD_DIR=${ROOT_DIR}/build
GO_FILES = $(wildcard *.go)
MODULE = $(shell go list -m)
VERSION = $(shell cat VERSION)
TAG := "v$(VERSION)"
ZIP_ARCHIVE := "$(TAG).zip"
TAR_GZ_ARCHIVE := "$(TAG).tar.gz"
RELEASE_NOTES = RELEASE_NOTES.md
PACKAGES := $(shell go list ./... | grep -v /vendor/)
LD_FLAGS := "-X ${NAME}/build.User=$(USER) -X main.Version=${VERSION} -X main.NAME=${NAME} -X main.BuildTime=`TZ=${TZ} date -u '+%Y-%m-%dT%H:%M:%SZ'`"
CUSTOM_FLAG := scan
FILES := $(shell git ls-files)
ASSETS := ${NAME}-v$(VERSION)-linux $(ZIP_ARCHIVE) $(TAR_GZ_ARCHIVE)

ifeq ($(strip $(ENV_FILE)),)
	CUSTOM_FLAG += -e .env
else
	CUSTOM_FLAG += -e $(ENV_FILE)
endif

default: help
.PHONY: default

# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Help information about make commands.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sed 's/Makefile://g' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
.PHONY: help

go-fmt: ## Formats all Go files in your project using the 'go fmt' tool.
	go fmt ./...
.PHONY: go-fmt

go-download: ## Downloads and installs the dependencies of your Go project using the 'go mod download' command.
	go mod download
.PHONY: go-download

go-get: ## Downloads and installs(update) the dependencies of your Go project using the 'go get' command.
	go get ./...
.PHONY: go-get

run: go-fmt go-get go-download ## Compiles and runs the Go program defined in main.go.
	@go run -ldflags ${LD_FLAGS} main.go ${CUSTOM_FLAG}
.PHONY: run

run-live: go-fmt go-get go-download ## Runs the Go program in live-reload mode using the 'air' tool with a configuration defined in ROOT_DIR/.air.toml.
	air -c ${ROOT_DIR}/.air.toml -- ${CUSTOM_FLAG} 
.PHONY: run-live

build-docker: ## Builds a Docker image with the name specified in the NAME variable, using the Dockerfile in the current directory.
	docker build -f Dockerfile -t $(NAME) .
.PHONY: build-docker

build: go-fmt go-get clean ## Formats the Go source code, downloads and installs the dependencies, cleans the build directory, and compiles the Go project for several different operating systems and architectures. The executables are placed in the build directory.
	mkdir -p ${BUILD_DIR}
	GOOS=linux GOARCH=amd64 go build -ldflags ${LD_FLAGS} -o ${BUILD_DIR}/${NAME}-v$(VERSION)-linux-amd64 ${ROOT_DIR}/*.go
	GOOS=linux GOARCH=arm64 go build -ldflags ${LD_FLAGS} -o ${BUILD_DIR}/${NAME}-v$(VERSION)-linux-arm64 ${ROOT_DIR}/*.go
	GOOS=windows GOARCH=amd64 go build -ldflags ${LD_FLAGS} -o ${BUILD_DIR}/${NAME}-v$(VERSION)-windows-amd64.exe ${ROOT_DIR}/*.go
	GOOS=darwin GOARCH=amd64 go build -ldflags ${LD_FLAGS} -o ${BUILD_DIR}/${NAME}-v$(VERSION)-macos-amd64 ${ROOT_DIR}/*.go
.PHONY: build

clean: ## Removes the build directory specified by the BUILD_DIR variable.
	rm -rf ${BUILD_DIR}
.PHONY: clean

version: ## Prints the current version of your project, which is stored in the VERSION variable.
	@echo $(VERSION)
.PHONY: version

version-major: ## Increments the major part of your project's version (e.g., 1.2.3 to 2.0.0) and updates the VERSION file.
	$(eval NEW_VERSION := $(shell echo "$(VERSION)" | awk -F. -v OFS=. '{$$1 = $$1 + 1; $$2 = 0; $$3 = 0; print $$0}'))
	@echo "Bumping version to $(NEW_VERSION)"
	$(eval VERSION := $(NEW_VERSION))
	@echo $(VERSION) > VERSION
.PHONY: version-major

version-minor: ## Increments the minor part of your project's version (e.g., 1.2.3 to 1.3.0) and updates the VERSION file.
	$(eval NEW_VERSION := $(shell echo "$(VERSION)" | awk -F. -v OFS=. '{$$2 = $$2 + 1; $$3 = 0; print $$0}'))
	@echo "Bumping version to $(NEW_VERSION)"
	$(eval VERSION := $(NEW_VERSION))
	@echo $(VERSION) > VERSION
.PHONY: version-minor

version-patch:  ## Increments the patch part of your project's version (e.g., 1.2.3 to 1.2.4) and updates the VERSION file.
	$(eval NEW_VERSION := $(shell echo "$(VERSION)" | awk -F. -v OFS=. '{$$3 = $$3 + 1; print $$0}'))
	@echo "Bumping version to $(NEW_VERSION)"
	$(eval VERSION := $(NEW_VERSION))
	@echo $(VERSION) > VERSION
.PHONY: version-patch

release: build ## Creates a new release on GitHub with your project's current version. It builds your project, creates the release, and uploads each file in the build directory to this release.
	$(eval RELEASE_EXISTS := $(shell gh release view v$(VERSION) > /dev/null 2>&1; echo $$?))
	@if [ $(RELEASE_EXISTS) -eq 0 ]; then \
		echo "Release v$(VERSION) already exists. Please update the version."; \
		exit 1; \
	fi
	@if [ ! -f $(RELEASE_NOTES) ]; then \
		echo "RELEASE_NOTES.md not found. Do you want to create it? (y/n)"; \
		read ans; \
		if [ $$ans = "y" ]; then \
			touch $(RELEASE_NOTES); \
		else \
			echo "Cannot proceed without RELEASE_NOTES.md"; \
			exit 1; \
		fi; \
		echo "The RELEASE_NOTES.md file has been created. Please fill it out before proceeding."; \
		exit 1; \
	fi
	gh release create v$(VERSION) --title "v$(VERSION)" --notes-file "$(RELEASE_NOTES)"
	for file in $(BUILD_DIR)/*; do \
    gh release upload v$(VERSION) "$$file"; \
  done
	rm $(RELEASE_NOTES)
.PHONY: release


zip: ## Creates a zip archive of the files specified by the FILES variable in the build directory.
	@cd build && zip $(ZIP_ARCHIVE) $(FILES)
.PHONY: zip

tar: ## Creates a tar.gz archive of the files specified by the FILES variable in the build directory.
	@cd build && tar -czvf $(TAR_GZ_ARCHIVE) $(FILES)
.PHONY: tar