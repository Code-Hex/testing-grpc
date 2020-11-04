package main

import (
	"context"
	"errors"
	"io"

	"github.com/Code-Hex/testing-grpc/internal/testing"
	"github.com/Songmu/prompter"
)

func (c *Client) runServerStream(ctx context.Context) error {
	strm, err := c.StreamClient.ServerStream(ctx, &testing.ServerStreamRequest{
		Duration: prompter.Prompt("duration", "1s"),
	})
	if err != nil {
		return err
	}
	for {
		ctx := strm.Context()
		select {
		case <-ctx.Done():
			defer strm.CloseSend()
			return ctx.Err()
		default:
		}
		var resp testing.ServerStreamResponse
		err := strm.RecvMsg(&resp)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		log.Info().Interface("response", &resp).Msg("success")
	}
	return nil
}
