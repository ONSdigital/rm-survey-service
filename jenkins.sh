# Set up environment
export PATH=$PATH:/usr/local/bin/go/bin
export PROJ=`git config --get remote.origin.url | sed 's/^https:\/\///' | sed 's/\.git$//' | tr '[:upper:]' '[:lower:]'`
export GOPATH=`pwd`

# Ensure a fully statically-linked binary that runs on different Linux distros
export CGO_ENABLED=0

# Clean directory
git clean -dfx

# Move project into GOPATH
mkdir -p src/$PROJ
ls -1 | grep -v ^src | xargs -I{} mv {} src/$PROJ/

# Download dependencies
go get github.com/tools/godep
export PATH=$PATH:$WORKSPACE/bin
cd src
godep get -v github.com/onsdigital/rm-survey-service

# Build the project
cd github.com/onsdigital/rm-survey-service
make build

# Create TAR file and upload to Artifactory
tar -cvf surveysvc.tar build/linux-amd64/bin/*
curl -u build:$PASSWORD -X PUT "http://artifactory.rmdev.onsdigital.uk/artifactory/libs-snapshot-local/uk/gov/ons/ctp/product/surveysvc/surveysvc-SNAPSHOT-$BUILD_NUMBER.tar" -T surveysvc.tar