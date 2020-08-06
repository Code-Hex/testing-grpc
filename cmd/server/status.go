package main

import (
	"context"

	"github.com/Code-Hex/testing-grpc/internal/testing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ testing.StatusServer = (*Status)(nil)

type Status struct{}

func (s *Status) Get(ctx context.Context, req *testing.StatusGetRequest) (*testing.StatusGetResponse, error) {
	if err := convertToCode(req.Code); err != nil {
		return nil, err
	}
	return &testing.StatusGetResponse{
		Msg: "Hello, World",
	}, nil
}

func codeErr(code codes.Code) error {
	return status.Error(code, code.String())
}

func convertToCode(c testing.StatusGetRequest_Code) error {
	switch c {
	case testing.StatusGetRequest_CANCELED:
		return codeErr(codes.Canceled)
	case testing.StatusGetRequest_UNKNOWN:
		return codeErr(codes.Unknown)
	case testing.StatusGetRequest_INVALIDARGUMENT:
		return codeErr(codes.InvalidArgument)
	case testing.StatusGetRequest_DEADLINE_EXCEEDED:
		return codeErr(codes.DeadlineExceeded)
	case testing.StatusGetRequest_NOT_FOUND:
		return codeErr(codes.NotFound)
	case testing.StatusGetRequest_ALREADY_EXISTS:
		return codeErr(codes.AlreadyExists)
	case testing.StatusGetRequest_PERMISSION_DENIED:
		return codeErr(codes.PermissionDenied)
	case testing.StatusGetRequest_RESOURCE_EXHAUSTED:
		return codeErr(codes.ResourceExhausted)
	case testing.StatusGetRequest_FAILED_PRECONDITION:
		return codeErr(codes.FailedPrecondition)
	case testing.StatusGetRequest_ABORTED:
		return codeErr(codes.Aborted)
	case testing.StatusGetRequest_OUT_OF_RANGE:
		return codeErr(codes.OutOfRange)
	case testing.StatusGetRequest_UNIMPLEMENTED:
		return codeErr(codes.Unimplemented)
	case testing.StatusGetRequest_INTERNAL:
		return codeErr(codes.Internal)
	case testing.StatusGetRequest_UNAVAILABLE:
		return codeErr(codes.Unavailable)
	case testing.StatusGetRequest_DATALOSS:
		return codeErr(codes.DataLoss)
	case testing.StatusGetRequest_UNAUTHENTICATED:
		return codeErr(codes.Unauthenticated)
	}
	return nil
}
