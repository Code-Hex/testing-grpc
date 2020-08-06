package main

import (
	"context"

	"github.com/Code-Hex/testing-grpc/internal/testing"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/pkg/errors"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var healthStatus = []testing.SetRequest_HealthCheckStatus{
	testing.SetRequest_UNKNOWN,
	testing.SetRequest_SERVING,
	testing.SetRequest_NOT_SERVING,
	testing.SetRequest_SERVICE_UNKNOWN,
}

func (c *Client) runChangeHealth(ctx context.Context) error {
	idx, err := fuzzyfinder.Find(healthStatus, func(i int) string {
		return healthStatus[i].String()
	})
	if err != nil {
		if errors.Is(err, fuzzyfinder.ErrAbort) {
			return nil
		}
		return errors.WithStack(err)
	}
	resp, err := c.ChangeHealth.Set(ctx, &testing.SetRequest{Status: healthStatus[idx]})
	if err != nil {
		loggingDetails(err)
	} else {
		log.Info().Interface("response", resp).Msg("success")
	}
	return nil
}

func (c *Client) runHealthClient(ctx context.Context) error {
	req := &healthpb.HealthCheckRequest{
		Service: testing.ChangeHealth,
	}
	resp, err := c.HealthClient.Check(ctx, req)
	if err != nil {
		loggingDetails(err)
	} else {
		log.Info().Stringer("response", resp.Status).Msg("success")
	}
	return nil
}
