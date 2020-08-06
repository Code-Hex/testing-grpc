package main

import (
	"context"

	"github.com/Code-Hex/testing-grpc/internal/testing"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/pkg/errors"
)

var details = []testing.DetailGetRequest_Code{
	testing.DetailGetRequest_OK,
	testing.DetailGetRequest_RETRY_INFO,
	testing.DetailGetRequest_DEBUG_INFO,
	testing.DetailGetRequest_QUOTA_FAILURE,
	testing.DetailGetRequest_ERROR_INFO,
	testing.DetailGetRequest_PRECONDITION_FAILURE,
	testing.DetailGetRequest_BAD_REQUEST,
	testing.DetailGetRequest_REQUEST_INFO,
	testing.DetailGetRequest_RESOURCE_INFO,
	testing.DetailGetRequest_HELP,
	testing.DetailGetRequest_LOCALIZED_MESSAGE,
	testing.DetailGetRequest_COMBINED_ALL,
}

func (c *Client) runDetailClient(ctx context.Context) error {
	idx, err := fuzzyfinder.Find(details, func(i int) string {
		return details[i].String()
	})
	if err != nil {
		if errors.Is(err, fuzzyfinder.ErrAbort) {
			return nil
		}
		return errors.WithStack(err)
	}
	req := &testing.DetailGetRequest{
		Code: details[idx],
	}
	resp, err := c.DetailClient.Get(ctx, req)
	if err != nil {
		loggingDetails(err)
	} else {
		log.Info().Interface("response", resp).Msg("success")
	}
	return nil
}
