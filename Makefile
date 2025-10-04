# Service version.
VERSION = 10.51.0

# Cross-compilation values.
ARCH=armd64
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

# This is the generic line to run all the tests bar those defined in vendor packages but the version
# of go used in cf (1.8) doesn't support multiple targets for test profile so each package will
# need to be explicitly enumerated for now (luckily there is only 1 ...)
# go test -race -coverprofile=coverage.txt -covermode=atomic $$(go list ./... | grep -v /vendor/)

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic github.com/ONSdigital/rm-survey-service/models

# Run integration and unit tests with the service running in docker.
# Builds a docker image, starts it up with a postgres container, waits for a successful response from the survey service
# /info endpoint, then runs unit tests and integration tests against the services in docker.
integration-test: docker
	docker compose -f compose-integration-tests.yml down
	docker compose -f compose-integration-tests.yml up -d
	./wait_for_startup_integration_tests.sh ||\
	(docker compose -f compose-integration-tests.yml down && exit 1)
	go test --tags=integration -race -coverprofile=coverage.txt -covermode=atomic github.com/ONSdigital/rm-survey-service/models ||\
	(docker compose -f compose-integration-tests.yml down && exit 1)
	docker compose -f compose-integration-tests.yml down

# Remove the build directory tree.
clean:
	if [ -d $(BUILD) ]; then rm -r $(BUILD); fi;

docker: build
	docker build . -t europe-west2-docker.pkg.dev/ons-ci-rmrasbs/images/survey
