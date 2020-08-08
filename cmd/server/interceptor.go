package main

import (
	"context"

	"github.com/Code-Hex/testing-grpc/internal/testing"
	"google.golang.org/grpc"
)

func DoLogServerInterceptor(s string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod != "/testing.Interceptor/Echo" {
			return handler(ctx, req)
		}
		log.Info().
			Interface("info", info).
			Interface("req", req).
			Msg(s)

		resp, err := handler(ctx, req)

		log.Info().
			Interface("response", resp).
			Err(err).
			Msg(s)

		return resp, err
	}
}

type Interceptor struct{}

var _ testing.InterceptorServer = (*Interceptor)(nil)

func (i *Interceptor) Echo(_ context.Context, req *testing.EchoRequest) (*testing.EchoResponse, error) {
	log.Info().
		Str("msg", req.Msg).
		Msg("called Echo()")
	return &testing.EchoResponse{Msg: req.Msg}, nil
}
