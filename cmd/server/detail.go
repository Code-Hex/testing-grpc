package main

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/Code-Hex/testing-grpc/internal/testing"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ testing.DetailServer = (*Detail)(nil)

type Detail struct{}

func (d *Detail) Get(ctx context.Context, req *testing.DetailGetRequest) (*testing.DetailGetResponse, error) {
	if err := convertToDetails(req.Code); err != nil {
		return nil, err
	}
	return &testing.DetailGetResponse{
		Msg: "OK details!!",
	}, nil
}

var details = map[testing.DetailGetRequest_Code]proto.Message{
	testing.DetailGetRequest_RETRY_INFO: &errdetails.RetryInfo{
		RetryDelay: &duration.Duration{
			Seconds: 2,
			Nanos:   100,
		},
	},
	testing.DetailGetRequest_DEBUG_INFO: &errdetails.DebugInfo{
		StackEntries: stackTraces(),
		Detail:       "debug info testing",
	},
	testing.DetailGetRequest_QUOTA_FAILURE: &errdetails.QuotaFailure{
		Violations: []*errdetails.QuotaFailure_Violation{
			{
				Subject:     "memory quota",
				Description: "used 1GB",
			},
			{
				Subject:     "API quota",
				Description: "used 100request/1month",
			},
		},
	},
	testing.DetailGetRequest_ERROR_INFO: &errdetails.ErrorInfo{
		Reason: "i/o timeout between application and database",
		Domain: "items",
		Metadata: map[string]string{
			"query":    "SELECT * FROM items",
			"function": "makeItem",
		},
	},
	testing.DetailGetRequest_PRECONDITION_FAILURE: &errdetails.PreconditionFailure{
		Violations: []*errdetails.PreconditionFailure_Violation{
			{
				Type:        "ENUM_USER_SERVICE_DOWN",
				Subject:     "user-service",
				Description: "Terms of service not accepted",
			},
		},
	},
	testing.DetailGetRequest_BAD_REQUEST: &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			{
				Field:       "request.item.id",
				Description: "invalid format",
			},
			{
				Field:       "request.category.sku_id",
				Description: "unexpected format: expected uppercases",
			},
		},
	},
	testing.DetailGetRequest_REQUEST_INFO: &errdetails.RequestInfo{
		RequestId:   "8DA1D58282DD43138804B7E75C86A50F",
		ServingData: "c3RhY2t0cmFjZQo=",
	},
	testing.DetailGetRequest_RESOURCE_INFO: &errdetails.ResourceInfo{
		ResourceType: "file",
		ResourceName: "codeowners",
		Owner:        "codehex",
		Description:  "permission denied",
	},
	testing.DetailGetRequest_HELP: &errdetails.Help{
		Links: []*errdetails.Help_Link{
			{
				Description: "please contact users team",
				Url:         "http://wiki.users-team.com",
			},
		},
	},
	testing.DetailGetRequest_LOCALIZED_MESSAGE: &errdetails.LocalizedMessage{
		Locale:  "en-US",
		Message: "message to en-US",
	},
}

func convertToDetails(c testing.DetailGetRequest_Code) error {
	switch c {
	case testing.DetailGetRequest_RETRY_INFO:
		return makeDetailsErr(codes.Unavailable, "retry please", details[c])
	case testing.DetailGetRequest_DEBUG_INFO:
		return makeDetailsErr(codes.Internal, "something wrong", details[c])
	case testing.DetailGetRequest_QUOTA_FAILURE:
		return makeDetailsErr(codes.Unavailable, "limit exceeded", details[c])
	case testing.DetailGetRequest_ERROR_INFO:
		return makeDetailsErr(codes.Internal, "caused internal error", details[c])
	case testing.DetailGetRequest_PRECONDITION_FAILURE:
		return makeDetailsErr(codes.FailedPrecondition, "caused some error", details[c])
	case testing.DetailGetRequest_BAD_REQUEST:
		return makeDetailsErr(codes.InvalidArgument, "invalid retuest fields", details[c])
	case testing.DetailGetRequest_REQUEST_INFO:
		return makeDetailsErr(codes.Internal, "something wrong", details[c])
	case testing.DetailGetRequest_RESOURCE_INFO:
		return makeDetailsErr(codes.PermissionDenied, "resource error", details[c])
	case testing.DetailGetRequest_HELP:
		return makeDetailsErr(codes.Unavailable, "temporary unavailable", details[c])
	case testing.DetailGetRequest_LOCALIZED_MESSAGE:
		return makeDetailsErr(codes.Internal, "something wrong", details[c])
	case testing.DetailGetRequest_COMBINED_ALL:
		keys := make([]testing.DetailGetRequest_Code, 0, len(details))
		for key := range details {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
		d := make([]proto.Message, len(details))
		for i, key := range keys {
			d[i] = details[key]
		}
		return makeDetailsErr(codes.Unknown, "combined all details", d...)
	}
	return nil
}

func stackTraces() []string {
	pc := make([]uintptr, 3)
	n := runtime.Callers(0, pc)
	if n == 0 {
		return []string{}
	}
	ret := make([]string, 0, 3)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		// To keep this example's output stable
		// even if there are changes in the testing package,
		// stop unwinding when we leave package runtime.
		if !strings.Contains(frame.File, "runtime/") {
			break
		}
		ret = append(ret,
			fmt.Sprintf(
				"file: %s, line: %d, %s",
				frame.File, frame.Line, frame.Function,
			),
		)
		if !more {
			break
		}
	}
	return ret
}

func makeDetailsErr(code codes.Code, msg string, d ...proto.Message) error {
	st, err := status.New(code, msg).WithDetails(d...)
	if err != nil {
		return errors.WithStack(err)
	}
	return st.Err()
}
