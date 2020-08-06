package main

import (
	"context"

	"github.com/Code-Hex/testing-grpc/internal/testing"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/pkg/errors"
)

var statuses = []testing.StatusGetRequest_Code{
	testing.StatusGetRequest_OK,
	testing.StatusGetRequest_CANCELED,
	testing.StatusGetRequest_UNKNOWN,
	testing.StatusGetRequest_INVALIDARGUMENT,
	testing.StatusGetRequest_DEADLINE_EXCEEDED,
	testing.StatusGetRequest_NOT_FOUND,
	testing.StatusGetRequest_ALREADY_EXISTS,
	testing.StatusGetRequest_PERMISSION_DENIED,
	testing.StatusGetRequest_RESOURCE_EXHAUSTED,
	testing.StatusGetRequest_FAILED_PRECONDITION,
	testing.StatusGetRequest_ABORTED,
	testing.StatusGetRequest_OUT_OF_RANGE,
	testing.StatusGetRequest_UNIMPLEMENTED,
	testing.StatusGetRequest_INTERNAL,
	testing.StatusGetRequest_UNAVAILABLE,
	testing.StatusGetRequest_DATALOSS,
	testing.StatusGetRequest_UNAUTHENTICATED,
}

func (c *Client) runStatusClient(ctx context.Context) error {
	idx, err := fuzzyfinder.Find(statuses, func(i int) string {
		return statuses[i].String()
	})
	if err != nil {
		if errors.Is(err, fuzzyfinder.ErrAbort) {
			return nil
		}
		return errors.WithStack(err)
	}
	req := &testing.StatusGetRequest{
		Code: statuses[idx],
	}
	resp, err := c.StatusClient.Get(ctx, req)
	if err != nil {
		loggingDetails(err)
	} else {
		log.Info().Interface("response", resp).Msg("success")
	}
	return nil
}
