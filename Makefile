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
	go install github.com/google/wire/cmd/wire@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/jaggerzhuang1994/kratos-foundation/cmd/protoc-gen-kratos-foundation-errors@main
	go install github.com/jaggerzhuang1994/kratos-foundation/cmd/protoc-gen-kratos-foundation-client@main
	go install github.com/jaggerzhuang1994/kratos-foundation/cmd/protoc-gen-jsonschema@main
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: config
# 生成配置
config:
	@echo "> protoc internal/conf & config.shcema.json ..."
	@if [ -n "$(INTERNAL_PROTO_FILES)" ]; then \
		protoc --proto_path=./internal \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./internal \
		   --jsonschema_out=./internal \
		   --jsonschema_opt=draft=Draft07 \
		   --jsonschema_opt=output_file_suffix=.schema.json \
		   --jsonschema_opt=preserve_proto_field_names=true \
		   --jsonschema_opt=merge=https://raw.githubusercontent.com/jaggerzhuang1994/kratos-foundation/refs/heads/main/config.schema.json \
	       $(INTERNAL_PROTO_FILES) && echo 'done'; \
	else \
		echo "no internal/conf proto files, skip"; \
	fi


.PHONY: api
# 生成api
api:
	@echo "> protoc api..."
	@if [ -n "$(API_PROTO_FILES)" ]; then \
		protoc --proto_path=./api \
			--proto_path=./third_party \
			--go_out=module=$(GO_MODULE_NAME)/api:./api \
			--go-http_out=module=$(GO_MODULE_NAME)/api:./api \
			--go-grpc_out=module=$(GO_MODULE_NAME)/api:./api \
			--validate_out=module=$(GO_MODULE_NAME)/api,lang=go:./api \
			--kratos-foundation-client_out=module=$(GO_MODULE_NAME)/api:./api \
			--kratos-foundation-errors_out=module=$(GO_MODULE_NAME)/api:./api \
			--openapiv2_out=$(OPEN_API_V2_FLAGS),merge_file_name=openapi.yaml:./ \
			$(API_PROTO_FILES) && echo 'done'; \
	else \
		echo "no api proto files, skip"; \
	fi

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

# 定义构建模板
define BUILD_TEMPLATE
.PHONY: build-$(1)

# 构建 $(1)
build-$(1):
	@echo "build $(1)..."
	@mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./cmd/$(1)
endef

# 查找所有构建目标
TARGETS=$(shell find cmd -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)
# 为每个目标生成构建命令
$(foreach target,$(TARGETS),$(eval $(call BUILD_TEMPLATE,$(target))))

.PHONY: test
# 运行单元测试
test:
	@echo "> running tests..."
	@go test -v -race -cover ./...

.PHONY: test-coverage
# 生成测试覆盖率报告
test-coverage:
	@echo "> generating coverage report..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

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
