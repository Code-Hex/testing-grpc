package main

import (
	"context"

	"github.com/Code-Hex/testing-grpc/internal/testing"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ChangeHealth struct {
	healthcheck *health.Server
}

func newChangeHealth(healthcheck *health.Server) *ChangeHealth {
	healthcheck.SetServingStatus(
		testing.ChangeHealth,
		healthpb.HealthCheckResponse_SERVING,
	)
	return &ChangeHealth{healthcheck: healthcheck}
}

var _ testing.ChangeHealthServer = (*ChangeHealth)(nil)

func (c *ChangeHealth) Set(ctx context.Context, req *testing.SetRequest) (*emptypb.Empty, error) {
	c.healthcheck.SetServingStatus(
		testing.ChangeHealth,
		convToServingStatus(req.Status),
	)
	return &emptypb.Empty{}, nil
}

func convToServingStatus(s testing.SetRequest_HealthCheckStatus) healthpb.HealthCheckResponse_ServingStatus {
	switch s {
	case testing.SetRequest_SERVING:
		return healthpb.HealthCheckResponse_SERVING
	case testing.SetRequest_NOT_SERVING:
		return healthpb.HealthCheckResponse_NOT_SERVING
	case testing.SetRequest_SERVICE_UNKNOWN:
		return healthpb.HealthCheckResponse_SERVICE_UNKNOWN
	}
	return healthpb.HealthCheckResponse_UNKNOWN
}
