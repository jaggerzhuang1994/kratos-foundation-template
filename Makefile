VERSION=$(shell git describe --tags --always)
#output_format=输出格式 enums_as_ints=枚举值输出为int omit_enum_default_value=是否忽略枚举默认值 allow_merge=合并文档
OPEN_API_V2_FLAGS=output_format=yaml,enums_as_ints=true,omit_enum_default_value=false,allow_merge=true
GO_MODULE_NAME=$(shell go list -m)

ifeq ($(shell go env GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: init
# 安装工具和环境
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/jaggerzhuang1994/kratos-foundation/cmd/protoc-gen-kratos-foundation-errors@latest
	go install github.com/jaggerzhuang1994/kratos-foundation/cmd/protoc-gen-kratos-foundation-client@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: config
# 生成配置
config:
	@echo "> protoc config..."
	@protoc --proto_path=./internal \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# 生成api
api:
	@echo "> protoc api..."
	@protoc --proto_path=./api \
			--proto_path=./third_party \
			--go_out=module=$(GO_MODULE_NAME)/api:./api \
			--go-http_out=module=$(GO_MODULE_NAME)/api:./api \
			--go-grpc_out=module=$(GO_MODULE_NAME)/api:./api \
			--validate_out=module=$(GO_MODULE_NAME)/api,lang=go:./api \
			--kratos-foundation-client_out=module=$(GO_MODULE_NAME)/api:./api \
			--kratos-foundation-errors_out=module=$(GO_MODULE_NAME)/api:./api \
			--openapiv2_out=$(OPEN_API_V2_FLAGS),merge_file_name=openapi.yaml:./ \
			$(API_PROTO_FILES)

.PHONY: generate
# 生成wire/其他生成
generate:
	@echo "> generate..."
	@go mod tidy
	@go generate ./...
	@go mod tidy

.PHONY: lint
# 代码审查
lint:
	@echo "> lint..."
	@golangci-lint run && echo 'lint ok'

.PHONY: all
# 生成所有/代码审查
all: api config generate lint

.PHONY: run
# 开发运行
run:
	@GOFLAGS="'-ldflags=-X=main.Version=$(VERSION)'" kratos run

.PHONY: build
# 构建
build:
	@echo "build..."
	@mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: help
# 展示帮助
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

# --- 通用 sed -i 封装（兼容 macOS / Linux / Git Bash）---
UNAME_S := $(shell uname -s)

SED_INPLACE := sed -i
ifeq ($(UNAME_S),Darwin)
	# macOS 的 sed 需要一个空备份后缀
	SED_INPLACE := sed -i ''
endif

.PHONY: rename-module
# 替换整个 go module 名称
# 用法：make rename-module NEW=github.com/jaggerzhuang1994/kratos-foundation-template
rename-module:
	@if [ -z "$(NEW)" ]; then \
		echo "用法: make rename-module NEW=github.com/jaggerzhuang1994/kratos-foundation-template"; \
		exit 1; \
	fi
	@old=$$(go list -m); \
	echo "old module: $$old"; \
	echo "new module: $(NEW)"; \
	echo ">> 更新 go.mod 中的 module 行"; \
	$(SED_INPLACE) "s#^module $$old#module $(NEW)#" go.mod; \
	echo ">> 全局替换源码中的导入路径等引用"; \
	find . -type f \( \
		-name '*.go' -o -name '*.proto' -o -name '*.yaml' -o -name '*.yml' -o -name 'Makefile' \
	\) \
		-not -path './.git/*' \
		-not -path './vendor/*' \
		-not -path './bin/*' \
		-print0 | xargs -0 $(SED_INPLACE) "s#$$old#$(NEW)#g"; \
	echo ">> make all..."; \
	make all
