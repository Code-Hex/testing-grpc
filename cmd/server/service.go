package main

import (
	"context"

	"github.com/Code-Hex/testing-grpc/internal/test"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ test.StatusServer = (*Status)(nil)

type Status struct{}

func (s *Status) Get(ctx context.Context, req *test.StatusGetRequest) (*test.StatusGetResponse, error) {
	if err := convertToCode(req.Code); err != nil {
		return nil, err
	}
	return &test.StatusGetResponse{
		Msg: "Hello, World",
	}, nil
}

func codeErr(code codes.Code) error {
	return status.Error(code, code.String())
}

func convertToCode(c test.StatusGetRequest_Code) error {
	switch c {
	case test.StatusGetRequest_CANCELED:
		return codeErr(codes.Canceled)
	case test.StatusGetRequest_UNKNOWN:
		return codeErr(codes.Unknown)
	case test.StatusGetRequest_INVALIDARGUMENT:
		return codeErr(codes.InvalidArgument)
	case test.StatusGetRequest_DEADLINE_EXCEEDED:
		return codeErr(codes.DeadlineExceeded)
	case test.StatusGetRequest_NOT_FOUND:
		return codeErr(codes.NotFound)
	case test.StatusGetRequest_ALREADY_EXISTS:
		return codeErr(codes.AlreadyExists)
	case test.StatusGetRequest_PERMISSION_DENIED:
		return codeErr(codes.PermissionDenied)
	case test.StatusGetRequest_RESOURCE_EXHAUSTED:
		return codeErr(codes.ResourceExhausted)
	case test.StatusGetRequest_FAILED_PRECONDITION:
		return codeErr(codes.FailedPrecondition)
	case test.StatusGetRequest_ABORTED:
		return codeErr(codes.Aborted)
	case test.StatusGetRequest_OUT_OF_RANGE:
		return codeErr(codes.OutOfRange)
	case test.StatusGetRequest_UNIMPLEMENTED:
		return codeErr(codes.Unimplemented)
	case test.StatusGetRequest_INTERNAL:
		return codeErr(codes.Internal)
	case test.StatusGetRequest_UNAVAILABLE:
		return codeErr(codes.Unavailable)
	case test.StatusGetRequest_DATALOSS:
		return codeErr(codes.DataLoss)
	case test.StatusGetRequest_UNAUTHENTICATED:
		return codeErr(codes.Unauthenticated)
	}
	return nil
}
