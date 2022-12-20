package cmd

import (
	"flag"
	"fmt"
	"io"
)

type grpcConfig struct {
	server string
	method string
	body   string
}

func HandleGrpc(w io.Writer, args []string) error {
	gc := grpcConfig{}

	fs := flag.NewFlagSet("grpc", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&gc.method, "method", "", "method to call")
	fs.StringVar(&gc.body, "body", "", "body of request")
	fs.Usage = func() {
		var usageString = `
grpc: A gRPC client.
 
grpc: <options> server`
		fmt.Fprintf(w, usageString)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}
	gc.server = fs.Arg(0)
	fmt.Fprintln(w, "Executing grpc command")
	return nil
}
