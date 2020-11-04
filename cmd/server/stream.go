package main

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Code-Hex/testing-grpc/internal/testing"
)

type Stream struct{}

var _ testing.StreamServer = (*Stream)(nil)

func (s *Stream) ServerStream(req *testing.ServerStreamRequest, strm testing.Stream_ServerStreamServer) error {
	d, err := time.ParseDuration(req.GetDuration())
	if err != nil {
		return err
	}
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for i := 0; i < 10; i++ {
		ctx := strm.Context()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
		msg := fmt.Sprintf("response: %d", i)
		err := strm.Send(&testing.ServerStreamResponse{
			Msg: msg,
		})
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
