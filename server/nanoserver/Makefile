APP_RELATIVE_PATH := $(shell a=`basename $$PWD` && cd .. && b=`basename $$PWD` && echo $$b/$$a)
APP_NAME := $(shell echo $(APP_RELATIVE_PATH) | sed -En "s/\//-/p")
REPO ?= registry.cn-hangzhou.aliyuncs.com/hayoujin
VERSION := $(shell git describe --tags --always)
#<repository>/<project-name>:<environment>-<version>
DOCKER_IMAGE := $(REPO)/$(APP_NAME):$(VERSION)

# 先安装vscode插件和WSL(ttps://learn.microsoft.com/en-us/windows/wsl/install)
cli:
	go install github.com/bufbuild/buf/cmd/buf@latest

build:
	mkdir -p bin && GOOS=linux go build -ldflags "-s -w -X main.Version=$(VERSION)" -o ./bin/ ./...

build-debug:
	mkdir -p bin && GOOS=linux go build -gcflags "all=-N -l" -o ./bin/ ./...

docker:
	docker buildx build -f build/Dockerfile -t $(DOCKER_IMAGE) .

dev-up:
	docker compose up -d

dev-stop:
	docker compose down

dev-clean:
	docker compose rm -sfv