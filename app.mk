GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
BRANCH=$(shell git branch --show-current)
THIRD_PARTY_PATH=third_party


INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)
#ifeq ($(GOHOSTOS), windows)
#	#the `find.exe` is different from `find` in bash/shell.
#	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
#	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
#	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
#	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
#	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
#	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
#else
#	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
#	API_PROTO_FILES=$(shell find api -name *.proto)
#endif

# generate internal proto
config:
	protoc --proto_path=./internal \
	       --proto_path=../../third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)


# generate api proto
api:
	cd ../../api/ && \
	buf generate

# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...
	cd bin && for f in *; do if [ ! "$${f#kva-*}" = "$$f" ]; then continue; fi; mv "$$f" "kva-$(VERSION)"; done

# build
clean:
	rm -rf ./bin

# generate
generate:
	go generate ./...
	go mod tidy

ent:
	ent generate \
				./internal/data/ent/schema \
				--template ./internal/data/ent/template \
				--feature privacy \
				--feature entql \
				--feature sql/modifier \
				--feature sql/upsert \
				--feature sql/lock


all:
	make api;
	make config;
	make generate;
	make ent;
	make wire;


run:
	make api;
	make config;
	make ent;
	go run ./cmd -conf ./configs

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
