# Service version.
VERSION = 10.47.0

# Cross-compilation values.
ARCH=amd64
OS_LINUX=linux
OS_MAC=darwin

# Output directory structures.
BUILD=build
LINUX_BUILD_ARCH=$(BUILD)/$(OS_LINUX)-$(ARCH)
MAC_BUILD_ARCH=$(BUILD)/$(OS_MAC)-$(ARCH)

# Flags to pass to the Go linker using the -ldflags="-X ..." option.
PACKAGE_PATH=github.com/onsdigital/rm-survey-service
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
	GOOS=$(OS_LINUX) GOARCH=$(ARCH) go build -o $(LINUX_BUILD_ARCH)/bin/surveysvc -ldflags="-X $(BRANCH_FLAG) -X $(BUILT_FLAG) -X $(COMMIT_FLAG) -X $(ORIGIN_FLAG) -X $(VERSION_FLAG)" surveysvc.go
	GOOS=$(OS_MAC) GOARCH=$(ARCH) go build -o $(MAC_BUILD_ARCH)/bin/surveysvc -ldflags="-X $(BRANCH_FLAG) -X $(BUILT_FLAG) -X $(COMMIT_FLAG) -X $(ORIGIN_FLAG) -X $(VERSION_FLAG)" surveysvc.go

# Run the tests.
test:
	go test -cover *.go

# Remove the build directory tree.
clean:
	if [ -d $(BUILD) ]; then rm -r $(BUILD); fi;

# Run a build then push to Cloud Foundry.
push-ci: build
	cf target -s ci
	cf push -f manifest-ci.yml
	
push-demo: build
	cf target -s demo
	cf push -f manifest-demo.yml

push-dev: build
	cf target -s dev
	cf push surveysvc-dev

push-int: build
	cf target -s int
	cf push -f manifest-int.yml

push-test: build
	cf target -s test
	cf push -f manifest-test.yml
