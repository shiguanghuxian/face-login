default:
	@echo 'Usage of make: [ build | linux_build | windows_build | clean ]'

build: 
	@go build -ldflags "-X main.VERSION=1.0.0 -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`' -X main.GIT_HASH=`git rev-parse HEAD` -s" -o ./build/faceserver ./

linux_build: 
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.VERSION=1.0.0 -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`' -X main.GIT_HASH=`git rev-parse HEAD` -s" -o ./build/faceserver ./

windows_build: 
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-X main.VERSION=1.0.0 -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`' -X main.GIT_HASH=`git rev-parse HEAD` -s" -o ./build/faceserver.exe ./

clean: 
	@rm -f ./build/faceserver*
	@rm -f ./build/logs/*.log

.PHONY: default build linux_build windows_build clean