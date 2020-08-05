package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Code-Hex/testing-grpc/internal/test"
	"github.com/Songmu/prompter"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

var statuses = []test.StatusGetRequest_Code{
	test.StatusGetRequest_OK,
	test.StatusGetRequest_CANCELED,
	test.StatusGetRequest_UNKNOWN,
	test.StatusGetRequest_INVALIDARGUMENT,
	test.StatusGetRequest_DEADLINE_EXCEEDED,
	test.StatusGetRequest_NOT_FOUND,
	test.StatusGetRequest_ALREADY_EXISTS,
	test.StatusGetRequest_PERMISSION_DENIED,
	test.StatusGetRequest_RESOURCE_EXHAUSTED,
	test.StatusGetRequest_FAILED_PRECONDITION,
	test.StatusGetRequest_ABORTED,
	test.StatusGetRequest_OUT_OF_RANGE,
	test.StatusGetRequest_UNIMPLEMENTED,
	test.StatusGetRequest_INTERNAL,
	test.StatusGetRequest_UNAVAILABLE,
	test.StatusGetRequest_DATALOSS,
	test.StatusGetRequest_UNAUTHENTICATED,
}

func run(ctx context.Context) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	conn, err := grpc.Dial(":"+port, grpc.WithInsecure())
	if err != nil {
		return errors.WithStack(err)
	}

	statusClient := test.NewStatusClient(conn)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)

LOOP:
	for {
		select {
		case <-sigCh:
			break LOOP
		default:
		}

		idx, err := fuzzyfinder.Find(statuses, func(i int) string {
			return statuses[i].String()
		})
		if err != nil {
			if errors.Is(err, fuzzyfinder.ErrAbort) {
				break LOOP
			}
			return errors.WithStack(err)
		}

		req := &test.StatusGetRequest{
			Code: statuses[idx],
		}
		resp, err := statusClient.Get(ctx, req)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(resp)
		}
		if !prompter.YN("continue?", true) {
			break
		}
	}

	return nil
}
