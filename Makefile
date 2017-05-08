BUILD=build
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)

export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)

build:
	@mkdir -p $(BUILD_ARCH)
	go build -o $(BUILD_ARCH)/bin/surveysvc survey-api/main.go

test:
	go test -cover survey-api/*.go

clean:
	test -d $(BUILD) && rm -r $(BUILD)