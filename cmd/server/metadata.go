package main

import (
	"context"
	"sort"

	"github.com/Code-Hex/testing-grpc/internal/testing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var _ testing.MetadataServer = (*Metadata)(nil)

type Metadata struct{}

func (m *Metadata) Get(ctx context.Context, _ *emptypb.Empty) (*testing.MetadataGetResponse, error) {
	// https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md#retrieving-metadata-from-context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.NotFound, "metadata is not found")
	}
	keys := make([]string, 0, len(md))
	for key := range md {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	entries := make([]*testing.MetadataGetResponse_Entry, len(keys))
	for i, key := range keys {
		entries[i] = &testing.MetadataGetResponse_Entry{
			Key:   key,
			Value: md[key],
		}
	}
	return &testing.MetadataGetResponse{
		Metadata: entries,
	}, nil
}
