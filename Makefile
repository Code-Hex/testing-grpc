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
	buf generate internal/testing

proto/clean:
	rm -f internal/testing/*.pb.go
