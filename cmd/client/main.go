package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Code-Hex/testing-grpc/internal/test"
	"github.com/Songmu/prompter"
	grpczerolog "github.com/cheapRoc/grpc-zerolog"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"

	// Necessary to print errdetails.
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
)

var log zerolog.Logger

func init() {
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	grpclog.SetLoggerV2(grpczerolog.New(log))
	fmt.Println()
}

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
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

	client := &Client{
		StatusClient: test.NewStatusClient(conn),
		DetailClient: test.NewDetailClient(conn),
	}
	reflectClient := newServerReflectionClient(ctx, conn)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		<-sigCh
		cancel()
	}()

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}

		// List gRPC services
		services, err := reflectClient.ListServices()
		if err != nil {
			if errors.Is(err, fuzzyfinder.ErrAbort) {
				return nil
			}
			return errors.WithStack(err)
		}

		// grpc.reflection.v1alpha.ServerReflection
		si, err := fuzzyfinder.Find(services, func(i int) string {
			return services[i]
		})
		if err != nil {
			return errors.WithStack(err)
		}

		switch services[si] {
		case "test.Status":
			if err := client.runStatusClient(ctx); err != nil {
				return err
			}
		case "test.Detail":
			if err := client.runDetailClient(ctx); err != nil {
				return err
			}
		default:
			continue LOOP
		}
		if !prompter.YN("continue?", true) {
			break
		}
	}

	return nil
}

type Client struct {
	StatusClient test.StatusClient
	DetailClient test.DetailClient
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

func (c *Client) runStatusClient(ctx context.Context) error {
	idx, err := fuzzyfinder.Find(statuses, func(i int) string {
		return statuses[i].String()
	})
	if err != nil {
		if errors.Is(err, fuzzyfinder.ErrAbort) {
			return nil
		}
		return errors.WithStack(err)
	}
	req := &test.StatusGetRequest{
		Code: statuses[idx],
	}
	resp, err := c.StatusClient.Get(ctx, req)
	if err != nil {
		loggingDetails(err)
	} else {
		log.Info().Interface("response", resp).Msg("success")
	}
	return nil
}

var details = []test.DetailGetRequest_Code{
	test.DetailGetRequest_OK,
	test.DetailGetRequest_RETRY_INFO,
	test.DetailGetRequest_DEBUG_INFO,
	test.DetailGetRequest_QUOTA_FAILURE,
	test.DetailGetRequest_ERROR_INFO,
	test.DetailGetRequest_PRECONDITION_FAILURE,
	test.DetailGetRequest_BAD_REQUEST,
	test.DetailGetRequest_REQUEST_INFO,
	test.DetailGetRequest_RESOURCE_INFO,
	test.DetailGetRequest_HELP,
	test.DetailGetRequest_LOCALIZED_MESSAGE,
	test.DetailGetRequest_COMBINED_ALL,
}

func (c *Client) runDetailClient(ctx context.Context) error {
	idx, err := fuzzyfinder.Find(details, func(i int) string {
		return details[i].String()
	})
	if err != nil {
		if errors.Is(err, fuzzyfinder.ErrAbort) {
			return nil
		}
		return errors.WithStack(err)
	}
	req := &test.DetailGetRequest{
		Code: details[idx],
	}
	resp, err := c.DetailClient.Get(ctx, req)
	if err != nil {
		loggingDetails(err)
	} else {
		log.Info().Interface("response", resp).Msg("success")
	}
	return nil
}

func newServerReflectionClient(ctx context.Context, conn *grpc.ClientConn) *grpcreflect.Client {
	cli := rpb.NewServerReflectionClient(conn)
	return grpcreflect.NewClient(ctx, cli)
}

func loggingDetails(err error) {
	logging := log.Error().Err(err)
	st, ok := status.FromError(err)
	if ok {
		for idx, d := range st.Details() {
			logging.Interface(
				fmt.Sprintf("details[%d]", idx),
				d,
			)
		}
	}
	logging.Msg("error")
}
