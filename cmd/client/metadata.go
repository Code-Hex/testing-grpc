package main

import (
	"context"
	"strings"

	"github.com/Code-Hex/testing-grpc/internal/testing"
	"github.com/Songmu/prompter"
	"google.golang.org/grpc/metadata"
)

func (c *Client) runMetadataClient(ctx context.Context) error {
	md := make([]string, 0)
	// https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md#sending-metadata
	for {
		key := prompter.Prompt("key", "default-key")
		value := prompter.Prompt("values (Become an array by comma-separated)", "default-value")
		vals := strings.Split(value, ",")
		for _, val := range vals {
			md = append(md, key, val)
		}
		if !prompter.YN("metadata continue?", true) {
			break
		}
	}
	ctx = metadata.AppendToOutgoingContext(ctx, md...)
	resp, err := c.MetadataClient.Get(ctx, &testing.MetadataGetRequest{})
	if err != nil {
		loggingDetails(err)
	} else {
		log.Info().Interface("response", resp).Msg("success")
	}
	return nil
}
