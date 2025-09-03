GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
BRANCH=$(shell git branch --show-current)
THIRD_PARTY_PATH=third_party

ifeq ($(GOHOSTOS), windows)
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
# initializer env
init:
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/google/gnostic@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install entgo.io/ent/cmd/ent@latest
	go mod tidy
	cd api/ && buf dep update

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	cd ../../api/ && \
	buf generate

#.PHONY: openapi
#openapi:
#	cd api/ && buf generate
#--template buf.kva.openapi.gen.yaml

.PHONY: deps
# update third_party proto
deps:
	buf export buf.build/bufbuild/protovalidate -o $(THIRD_PARTY_PATH)
	buf export buf.build/protocolbuffers/wellknowntypes -o $(THIRD_PARTY_PATH)
	buf export buf.build/googleapis/googleapis -o $(THIRD_PARTY_PATH)
	buf export buf.build/envoyproxy/protoc-gen-validate -o $(THIRD_PARTY_PATH)
	buf export buf.build/gnostic/gnostic -o $(THIRD_PARTY_PATH)
	buf export buf.build/kratos/apis -o $(THIRD_PARTY_PATH)

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...
	cd bin && for f in *; do if [ ! "$${f#kva-*}" = "$$f" ]; then continue; fi; mv "$$f" "kva-$(VERSION)"; done

.PHONY: clean
# build
clean:
	rm -rf ./bin


.PHONY: generate
# generate
generate:
	go generate ./...
	go mod tidy

.PHONY: all
# generate all
all:
	make api;
	make config;
	make generate;
	make wire;

.PHONY: run
# kratos run
run:
	kratos run

.PHONY: wire
wire:
	wire ./...


# show help
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
