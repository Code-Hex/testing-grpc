package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Code-Hex/testing-grpc/internal/testing"
	"github.com/Songmu/prompter"
	grpczerolog "github.com/cheapRoc/grpc-zerolog"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
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
		StatusClient:   testing.NewStatusClient(conn),
		DetailClient:   testing.NewDetailClient(conn),
		MetadataClient: testing.NewMetadataClient(conn),
		ChangeHealth:   testing.NewChangeHealthClient(conn),
		HealthClient:   healthpb.NewHealthClient(conn),
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
			return errors.WithStack(err)
		}

		// grpc.reflection.v1alpha.ServerReflection
		si, err := fuzzyfinder.Find(services, func(i int) string {
			return services[i]
		})
		if err != nil {
			if errors.Is(err, fuzzyfinder.ErrAbort) {
				return nil
			}
			return errors.WithStack(err)
		}

		switch services[si] {
		case testing.Status:
			if err := client.runStatusClient(ctx); err != nil {
				return err
			}
		case testing.Detail:
			if err := client.runDetailClient(ctx); err != nil {
				return err
			}
		case testing.Metadata:
			if err := client.runMetadataClient(ctx); err != nil {
				return err
			}
		case testing.ChangeHealth:
			if err := client.runChangeHealth(ctx); err != nil {
				return err
			}
		case "grpc.health.v1.Health":
			if err := client.runHealthClient(ctx); err != nil {
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
	StatusClient   testing.StatusClient
	DetailClient   testing.DetailClient
	MetadataClient testing.MetadataClient
	ChangeHealth   testing.ChangeHealthClient
	HealthClient   healthpb.HealthClient
}

var statuses = []testing.StatusGetRequest_Code{
	testing.StatusGetRequest_OK,
	testing.StatusGetRequest_CANCELED,
	testing.StatusGetRequest_UNKNOWN,
	testing.StatusGetRequest_INVALIDARGUMENT,
	testing.StatusGetRequest_DEADLINE_EXCEEDED,
	testing.StatusGetRequest_NOT_FOUND,
	testing.StatusGetRequest_ALREADY_EXISTS,
	testing.StatusGetRequest_PERMISSION_DENIED,
	testing.StatusGetRequest_RESOURCE_EXHAUSTED,
	testing.StatusGetRequest_FAILED_PRECONDITION,
	testing.StatusGetRequest_ABORTED,
	testing.StatusGetRequest_OUT_OF_RANGE,
	testing.StatusGetRequest_UNIMPLEMENTED,
	testing.StatusGetRequest_INTERNAL,
	testing.StatusGetRequest_UNAVAILABLE,
	testing.StatusGetRequest_DATALOSS,
	testing.StatusGetRequest_UNAUTHENTICATED,
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
	req := &testing.StatusGetRequest{
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

var details = []testing.DetailGetRequest_Code{
	testing.DetailGetRequest_OK,
	testing.DetailGetRequest_RETRY_INFO,
	testing.DetailGetRequest_DEBUG_INFO,
	testing.DetailGetRequest_QUOTA_FAILURE,
	testing.DetailGetRequest_ERROR_INFO,
	testing.DetailGetRequest_PRECONDITION_FAILURE,
	testing.DetailGetRequest_BAD_REQUEST,
	testing.DetailGetRequest_REQUEST_INFO,
	testing.DetailGetRequest_RESOURCE_INFO,
	testing.DetailGetRequest_HELP,
	testing.DetailGetRequest_LOCALIZED_MESSAGE,
	testing.DetailGetRequest_COMBINED_ALL,
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
	req := &testing.DetailGetRequest{
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

var healthStatus = []testing.SetRequest_HealthCheckStatus{
	testing.SetRequest_UNKNOWN,
	testing.SetRequest_SERVING,
	testing.SetRequest_NOT_SERVING,
	testing.SetRequest_SERVICE_UNKNOWN,
}

func (c *Client) runChangeHealth(ctx context.Context) error {
	idx, err := fuzzyfinder.Find(healthStatus, func(i int) string {
		return healthStatus[i].String()
	})
	if err != nil {
		if errors.Is(err, fuzzyfinder.ErrAbort) {
			return nil
		}
		return errors.WithStack(err)
	}
	resp, err := c.ChangeHealth.Set(ctx, &testing.SetRequest{Status: healthStatus[idx]})
	if err != nil {
		loggingDetails(err)
	} else {
		log.Info().Interface("response", resp).Msg("success")
	}
	return nil
}

func (c *Client) runHealthClient(ctx context.Context) error {
	req := &healthpb.HealthCheckRequest{
		Service: testing.ChangeHealth,
	}
	resp, err := c.HealthClient.Check(ctx, req)
	if err != nil {
		loggingDetails(err)
	} else {
		log.Info().Stringer("response", resp.Status).Msg("success")
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
