package cmd

import (
	"flag"
	"fmt"
	"io"
)

var methods = map[string]bool{
	"POST": true,
	"GET":  true,
	"HEAD": true,
}

type httpConfig struct {
	url  string
	verb string
}

func HandleHttp(w io.Writer, args []string) error {
	hc := httpConfig{}
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&hc.verb, "verb", "GET", "HTTP method")

	fs.Usage = func() {
		var usageString = `
http: A HTTP client.
 
http: <options> server`
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

	if !hc.validateMethod() {
		return InvalidHttpMethod
	}
	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}

	hc.url = fs.Arg(0)
	fmt.Fprintln(w, "Executing http command")
	return nil
}

func (c *httpConfig) validateMethod() bool {
	_, ok := methods[c.verb]

	return ok
}
