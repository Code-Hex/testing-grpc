package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func newServerReflectionClient(ctx context.Context, conn *grpc.ClientConn) *grpcreflect.Client {
	cli := rpb.NewServerReflectionClient(conn)
	return grpcreflect.NewClient(ctx, cli)
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
		grpc.WithBlock(),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	reflectClient := newServerReflectionClient(ctx, conn)

	// List gRPC services
	services, err := reflectClient.ListServices()
	if err != nil {
		return errors.WithStack(err)
	}
	svc, err := fuzzyFind(services, func(i int) string {
		return services[i]
	})
	if err != nil {
		return errors.WithStack(err)
	}
	sd, err := reflectClient.ResolveService(svc)
	if err != nil {
		return errors.WithStack(err)
	}

	methodDescs := sd.GetMethods()
	methods := make([]string, len(methodDescs))
	for i, v := range methodDescs {
		methods[i] = v.GetName()
	}

	mi, err := fuzzyfinder.Find(methods, func(i int) string {
		return methods[i]
	})
	if err != nil {
		return errors.WithStack(err)
	}
	md := methodDescs[mi]
	fmt.Println(md.GetFullyQualifiedName())

	reflSource := grpcurl.DescriptorSourceFromServer(ctx, reflectClient)
	headers := []string{}
	in := os.Stdin
	opts := grpcurl.FormatOptions{}
	rf, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.Format("json"), reflSource, in, opts)
	if err != nil {
		return errors.WithStack(err)
	}

	var buf bytes.Buffer
	h := grpcurl.NewDefaultEventHandler(
		&buf,
		reflSource,
		formatter,
		false, // verbose
	)
	if err := grpcurl.InvokeRPC(ctx, reflSource, conn, md.GetFullyQualifiedName(), headers, h, rf.Next); err != nil {
		return errors.WithStack(err)
	}

	if _, err := io.Copy(os.Stdout, &buf); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func fuzzyFind(s []string, f func(i int) string) (string, error) {
	i, err := fuzzyfinder.Find(s, f)
	if err != nil {
		// if errors.Is(err, fuzzyfinder.ErrAbort) {
		// 	return "", nil
		// }
		return "", err
	}
	return s[i], nil
}
