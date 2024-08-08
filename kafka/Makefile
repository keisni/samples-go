gomod_on = GO111MODULE=on
go_build_cmd = go build
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
build_path := $(dir $(mkfile_path))

BIN_DIR=$(build_path)/bin
GO_BUILD_LINUX_AMD64=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(gomod_on) $(go_build_cmd)
GO_BUILD_TAGS ?= -gcflags="all=-N -l"

init:
	mkdir -p ${BIN_DIR}

producer-starter:
	$(GO_BUILD_LINUX_AMD64) \
	-tags $(GO_BUILD_TAGS) -o ${BIN_DIR}/producer-starter \
	-installsuffix cgo $(build_path)/producer/starter

producer-worker:
	$(GO_BUILD_LINUX_AMD64) \
	-tags $(GO_BUILD_TAGS) -o ${BIN_DIR}/producer-worker \
	-installsuffix cgo $(build_path)/producer/worker

consumer-starter:
	$(GO_BUILD_LINUX_AMD64) \
	-tags $(GO_BUILD_TAGS) -o ${BIN_DIR}/consumer-starter \
	-installsuffix cgo $(build_path)/consumer/starter

consumer-worker:
	$(GO_BUILD_LINUX_AMD64) \
	-tags $(GO_BUILD_TAGS) -o ${BIN_DIR}/consumer-worker \
	-installsuffix cgo $(build_path)/consumer/worker

all: init producer-starter producer-worker consumer-starter consumer-worker