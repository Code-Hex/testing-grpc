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

proto/compile:
	mkdir -p $(OUTPUT_DIR)
	protoc -I. --go_out=plugins=grpc:$(OUTPUT_DIR) internal/test/*.proto
	cp $(OUTPUT_DIR)/$(PKG)/internal/test/*.go internal/test/

proto/clean:
	rm -f internal/testdata/*.pb.go
