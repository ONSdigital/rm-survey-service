# Service version.
VERSION = 10.47.2

# Cross-compilation values.
ARCH=amd64
OS_LINUX=linux
OS_MAC=darwin

# Output directory structures.
BUILD=build
LINUX_BUILD_ARCH=$(BUILD)/$(OS_LINUX)-$(ARCH)
MAC_BUILD_ARCH=$(BUILD)/$(OS_MAC)-$(ARCH)

# Flags to pass to the Go linker using the -ldflags="-X ..." option.
PACKAGE_PATH=github.com/ONSdigital/rm-survey-service
BRANCH_FLAG=$(PACKAGE_PATH)/models.branch=$(BRANCH)
BUILT_FLAG=$(PACKAGE_PATH)/models.built=$(BUILT)
COMMIT_FLAG=$(PACKAGE_PATH)/models.commit=$(COMMIT)
ORIGIN_FLAG=$(PACKAGE_PATH)/models.origin=$(ORIGIN)
VERSION_FLAG=$(PACKAGE_PATH)/models.version=$(VERSION)

# Get the Git branch the commit is from, stripping the leading asterisk.
export BRANCH?=$(shell git branch --contains $(COMMIT) | grep \* | cut -d ' ' -f2)

# Get the current date/time in UTC and ISO-8601 format.
export BUILT?=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# Get the full Git commit SHA-1 hash.
export COMMIT?=$(shell git rev-parse HEAD)

# Get the Git repo origin.
export ORIGIN?=$(shell git config --get remote.origin.url)

# Cross-compile the binary for Linux and macOS, setting linker flags for information returned by the GET /info endpoint.
build: clean
	GOOS=$(OS_LINUX) GOARCH=$(ARCH) go build -o $(LINUX_BUILD_ARCH)/bin/main -ldflags="-X $(BRANCH_FLAG) -X $(BUILT_FLAG) -X $(COMMIT_FLAG) -X $(ORIGIN_FLAG) -X $(VERSION_FLAG)" main.go
	GOOS=$(OS_MAC) GOARCH=$(ARCH) go build -o $(MAC_BUILD_ARCH)/bin/main -ldflags="-X $(BRANCH_FLAG) -X $(BUILT_FLAG) -X $(COMMIT_FLAG) -X $(ORIGIN_FLAG) -X $(VERSION_FLAG)" main.go

# Run the tests.
test:
	go test -race -coverprofile=coverage.txt -covermode=atomic $$(go list ./... | grep -v /vendor/)

# Remove the build directory tree.
clean:
	if [ -d $(BUILD) ]; then rm -r $(BUILD); fi;

docker: build
	docker build . -t sdcplatform/surveysvc
