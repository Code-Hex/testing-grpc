# testing-grpc

A server and client developed to understand the behavior of gRPC. Mainly intended to be useful for application development using gRPC.

[![asciicast](https://asciinema.org/a/351916.svg)](https://asciinema.org/a/351916)

## build

```
$ make build

# startup server (default port 3000)
$ ./bin/server

# startup client (default port 3000)
$ ./bin/client
```

If you want to change port, you can change the environment variable of `PORT`. and you can use `.env` file :D

## supported

- Unary
  - [x] status
  - [x] error details
  - [x] metadata
  - [x] [health check](https://github.com/grpc/grpc-go/tree/master/examples/features/health)
  - [ ] cancel, timeout
  - [ ] interceptor
  - [x] stats
    - If you want to disable logging, please set `LOG_LEVEL` to `INFO`. see `.env` file.

## gRPC documents

- [Error Handling – gRPC](https://grpc.io/docs/guides/error/)
- [go-grpc Learn More](https://github.com/grpc/grpc-go#learn-more)
  - [Low-level technical docs](https://github.com/grpc/grpc-go/tree/master/Documentation)
- [Language Guide (proto3)](https://developers.google.com/protocol-buffers/docs/proto3)
