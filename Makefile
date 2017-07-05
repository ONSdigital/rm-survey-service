# Output directory structure.
BUILD=build
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)

# Flags to pass to the Go linker using the -ldflags="-X ..." option.
PACKAGE_PATH = github.com/onsdigital/rm-survey-service
BRANCH_FLAG = $(PACKAGE_PATH)/models.branch=$(BRANCH)
BUILT_FLAG = $(PACKAGE_PATH)/models.built=$(BUILT)
COMMIT_FLAG = $(PACKAGE_PATH)/models.commit=$(COMMIT)
ORIGIN_FLAG = $(PACKAGE_PATH)/models.origin=$(ORIGIN)

# Get the operating system details as reported by the Go environment.
export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)

# Get the Git branch the commit is from, stripping the leading asterisk.
export BRANCH?=$(shell git branch --contains $(COMMIT) | grep \* | cut -d ' ' -f2)

# Get the current date/time in UTC and ISO-8601 format.
export BUILT?=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# Get the full Git commit SHA-1 hash.
export COMMIT?=$(shell git rev-parse HEAD)

# Get the Git repo origin.
export ORIGIN?=$(shell git remote get-url origin)

# Build the binary, setting linker flags for information returned by the GET /about endpoint.
build:	
	go build -o $(BUILD_ARCH)/bin/surveysvc -ldflags="-X $(BUILT_FLAG) -X $(COMMIT_FLAG) -X $(BRANCH_FLAG) -X $(ORIGIN_FLAG)" surveysvc.go

# Run the tests.
test:
	go test -cover *.go

# Remote the build directory tree.
clean:
	test -d $(BUILD) && rm -r $(BUILD)

# Update dependencies then push to Cloud Foundry.
push:
	godep get; godep save
	cf push
