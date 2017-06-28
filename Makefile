BUILD=build
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)

export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)

build:
	go build -o $(BUILD_ARCH)/bin/surveysvc surveysvc.go

test:
	go test -cover *.go

clean:
	test -d $(BUILD) && rm -r $(BUILD)

push:
	godep get; godep save
	cf push
