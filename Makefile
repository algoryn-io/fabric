# Algoryn Fabric — Protocol Buffers code generation
#
# Requires: protoc (https://grpc.io/docs/protoc-installation/)
#   macOS: brew install protobuf
#   Debian/Ubuntu: apt install protobuf-compiler
#
# Go plugin:
#   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# Ensure $(go env GOPATH)/bin is on PATH.

.PHONY: proto proto-check

PROTOC        ?= protoc
PROTO_ROOT    := proto
GO_OUT        := gen/go
PROTO_FILES   := $(wildcard $(PROTO_ROOT)/fabric/v1/*.proto)

# Well-known protos are vendored under proto/third_party/google/protobuf
# so generation does not depend on a system protobuf include tree.

.PHONY: proto proto-check

PROTOC        ?= protoc
PROTO_ROOT    := proto
GO_OUT        := gen/go

proto-check:
	@command -v $(PROTOC) >/dev/null 2>&1 || (echo "error: protoc not found; install from https://grpc.io/docs/protoc-installation/" && exit 1)
	@command -v protoc-gen-go >/dev/null 2>&1 || (echo "error: protoc-gen-go not found; run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest" && exit 1)

# Run from PROTO_ROOT so paths=source_relative emits gen/go/fabric/v1/*.pb.go
proto: proto-check
	@mkdir -p $(GO_OUT)
	cd $(PROTO_ROOT) && $(PROTOC) \
		-I=. \
		-I=third_party \
		--go_out=../$(GO_OUT) \
		--go_opt=paths=source_relative \
		fabric/v1/types_common.proto \
		fabric/v1/metrics.proto \
		fabric/v1/events.proto
