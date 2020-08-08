package main

import (
	"context"

	"github.com/Code-Hex/testing-grpc/internal/testing"
	"github.com/Songmu/prompter"
	"google.golang.org/grpc"
)

func DoLogServerInterceptor(s string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if method != "/testing.Interceptor/Echo" {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		log.Info().
			Interface("method", method).
			Interface("req", req).
			Msg(s)

		err := invoker(ctx, method, req, reply, cc, opts...)

		log.Info().
			Interface("reply", reply).
			Err(err).
			Msg(s)

		return err
	}
}

func (c *Client) runInterceptorClient(ctx context.Context) {
	msg := prompter.Prompt("message", "default-message")
	resp, err := c.InterceptorClient.Echo(ctx, &testing.EchoRequest{Msg: msg})
	if err != nil {
		loggingDetails(err)
	} else {
		log.Info().Interface("response", resp).Msg("success")
	}
}
