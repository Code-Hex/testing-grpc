package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Code-Hex/testing-grpc/internal/level"
	"github.com/Code-Hex/testing-grpc/internal/stats"
	"github.com/Code-Hex/testing-grpc/internal/testing"
	grpczerolog "github.com/cheapRoc/grpc-zerolog"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
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
	srv := newServer()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	ln, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Info().Str("port", port).Msg("Running server")

	go srv.Serve(ln)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
		log.Info().Msg("received SIGTERM, exiting server gracefully")
	case <-ctx.Done():
	}
	go srv.GracefulStop()

	return nil
}

func newServer() *grpc.Server {
	srv := grpc.NewServer(
		grpc.StatsHandler(stats.NewHandler(log)),
		grpc.ChainUnaryInterceptor(
			DoLogServerInterceptor("A"),
			DoLogServerInterceptor("B"),
		),
	)

	// health check service
	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(srv, healthcheck)

	// register gRPC services
	testing.RegisterStatusServer(srv, &Status{})
	testing.RegisterDetailServer(srv, &Detail{})
	testing.RegisterMetadataServer(srv, &Metadata{})
	testing.RegisterChangeHealthServer(srv, newChangeHealth(healthcheck))
	testing.RegisterInterceptorServer(srv, &Interceptor{})
	testing.RegisterStreamServer(srv, &Stream{})

	// enable reflection
	reflection.Register(srv)

	return srv
}
