package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Code-Hex/testing-grpc/internal/test"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return errors.WithStack(err)
	}

	go srv.Serve(ln)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
		log.Println("received SIGTERM, exiting server gracefully")
	case <-ctx.Done():
	}
	go srv.GracefulStop()

	return nil
}

func newServer() *grpc.Server {
	srv := grpc.NewServer()
	// enable reflection
	reflection.Register(srv)
	test.RegisterStatusServer(srv, &Status{})
	test.RegisterDetailServer(srv, &Detail{})
	return srv
}
