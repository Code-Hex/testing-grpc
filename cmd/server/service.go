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

func convertToCode(c test.Code) error {
	switch c {
	case test.Code_CANCELED:
		return codeErr(codes.Canceled)
	case test.Code_UNKNOWN:
		return codeErr(codes.Unknown)
	case test.Code_INVALIDARGUMENT:
		return codeErr(codes.InvalidArgument)
	case test.Code_DEADLINE_EXCEEDED:
		return codeErr(codes.DeadlineExceeded)
	case test.Code_NOT_FOUND:
		return codeErr(codes.NotFound)
	case test.Code_ALREADY_EXISTS:
		return codeErr(codes.AlreadyExists)
	case test.Code_PERMISSION_DENIED:
		return codeErr(codes.PermissionDenied)
	case test.Code_RESOURCE_EXHAUSTED:
		return codeErr(codes.ResourceExhausted)
	case test.Code_FAILED_PRECONDITION:
		return codeErr(codes.FailedPrecondition)
	case test.Code_ABORTED:
		return codeErr(codes.Aborted)
	case test.Code_OUT_OF_RANGE:
		return codeErr(codes.OutOfRange)
	case test.Code_UNIMPLEMENTED:
		return codeErr(codes.Unimplemented)
	case test.Code_INTERNAL:
		return codeErr(codes.Internal)
	case test.Code_UNAVAILABLE:
		return codeErr(codes.Unavailable)
	case test.Code_DATALOSS:
		return codeErr(codes.DataLoss)
	case test.Code_UNAUTHENTICATED:
		return codeErr(codes.Unauthenticated)
	}
	return nil
}
