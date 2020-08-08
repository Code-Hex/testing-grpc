package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Code-Hex/testing-grpc/internal/level"
	"github.com/Code-Hex/testing-grpc/internal/stats"
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
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"

	// Necessary to print errdetails.
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
)

var log zerolog.Logger

func init() {
	log = zerolog.New(
		zerolog.ConsoleWriter{
			Out: os.Stderr,
		}).
		With().
		Timestamp().
		Logger().
		Level(level.Log(os.Getenv("LOG_LEVEL")))
	grpclog.SetLoggerV2(grpczerolog.New(log))
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

	conn, err := grpc.Dial(":"+port,
		grpc.WithInsecure(),
		grpc.WithStatsHandler(stats.NewHandler(log)),
		grpc.WithChainUnaryInterceptor(
			DoLogServerInterceptor("A"),
			DoLogServerInterceptor("B"),
		),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	client := &Client{
		StatusClient:      testing.NewStatusClient(conn),
		DetailClient:      testing.NewDetailClient(conn),
		MetadataClient:    testing.NewMetadataClient(conn),
		ChangeHealth:      testing.NewChangeHealthClient(conn),
		HealthClient:      healthpb.NewHealthClient(conn),
		InterceptorClient: testing.NewInterceptorClient(conn),
	}
	reflectClient := newServerReflectionClient(ctx, conn)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		<-sigCh
		cancel()
	}()

	// List gRPC services
	services, err := reflectClient.ListServices()
	if err != nil {
		return errors.WithStack(err)
	}

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		default:
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
		case testing.Interceptor:
			client.runInterceptorClient(ctx)
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
	StatusClient      testing.StatusClient
	DetailClient      testing.DetailClient
	MetadataClient    testing.MetadataClient
	ChangeHealth      testing.ChangeHealthClient
	HealthClient      healthpb.HealthClient
	InterceptorClient testing.InterceptorClient
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
