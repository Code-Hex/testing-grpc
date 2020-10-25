PKG=github.com/Code-Hex/testing-grpc
OUTPUT_DIR=_output

.PHONY: build
build:
	@echo "+ $@"
	@echo "+ build server"
	go build -o bin/server -trimpath -mod=readonly \
        github.com/Code-Hex/testing-grpc/cmd/server
	@echo "+ build client"
	go build -o bin/client -trimpath -mod=readonly \
        github.com/Code-Hex/testing-grpc/cmd/client

.PHONY: build/reflect_cli
build/reflect_cli:
	@echo "+ build reflection client"
	go build -o bin/reflect_cli -trimpath -mod=readonly \
        github.com/Code-Hex/testing-grpc/cmd/reflect_cli

proto/compile:
	mkdir -p $(OUTPUT_DIR)
	protoc -I. -Ithird_party/protocolbuffers/src/google/protobuf --go_out=plugins=grpc:$(OUTPUT_DIR) internal/testing/*.proto
	cp $(OUTPUT_DIR)/$(PKG)/internal/testing/*.go internal/testing/

proto/clean:
	rm -f internal/testing/*.pb.go
