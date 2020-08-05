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

var statuses = []test.Code{
	test.Code_OK,
	test.Code_CANCELED,
	test.Code_UNKNOWN,
	test.Code_INVALIDARGUMENT,
	test.Code_DEADLINE_EXCEEDED,
	test.Code_NOT_FOUND,
	test.Code_ALREADY_EXISTS,
	test.Code_PERMISSION_DENIED,
	test.Code_RESOURCE_EXHAUSTED,
	test.Code_FAILED_PRECONDITION,
	test.Code_ABORTED,
	test.Code_OUT_OF_RANGE,
	test.Code_UNIMPLEMENTED,
	test.Code_INTERNAL,
	test.Code_UNAVAILABLE,
	test.Code_DATALOSS,
	test.Code_UNAUTHENTICATED,
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
