BUILD=build
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)

export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)

build:
	mkdir -p $(BUILD_ARCH)/bin/sql
	cp sql/bootstrap.sql $(BUILD_ARCH)/bin/sql
	go build -o $(BUILD_ARCH)/bin/surveysvc survey-api/main.go

test:
	go test -cover survey-api/*.go

clean:
	test -d $(BUILD) && rm -r $(BUILD)

push:
	cp -fr ./sql ./survey-api/sql
	cd survey-api; godep get; godep save
	cf push
	rm -rf ./survey-api/sql
