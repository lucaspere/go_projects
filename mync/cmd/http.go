package cmd

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var methods = map[string]bool{
	"POST": true,
	"GET":  true,
	"HEAD": true,
}

type httpConfig struct {
	url      string
	verb     string
	filePath string
}

func HandleHttp(w io.Writer, args []string) error {
	hc := httpConfig{}
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&hc.verb, "verb", "GET", "HTTP method")
	fs.StringVar(&hc.filePath, "output", "", "File path where save the data's output")
	fs.StringVar(&hc.filePath, "body", "", "JSON data to be used as payload")
	fs.StringVar(&hc.filePath, "body-file", "", "File path containing the JSON data to be used as payload")

	fs.Usage = func() {
		var usageString = `
http: A HTTP client.
http: <options> server`

		fmt.Fprintln(w, usageString)
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

	if len(hc.filePath) > 0 {
		var file *os.File

		path := filepath.Join(hc.filePath)
		file, err = os.Open(path)
		if err != nil {
			file, err = os.Create(path)
			if err != nil {
				return err
			}
		}

		hc.fetchData(file)
	} else {
		hc.fetchData(w)
	}

	return nil
}

func (c *httpConfig) fetchData(w io.Writer) error {
	req, err := http.NewRequest(c.verb, c.url, nil)
	if err != nil {
		return err
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	w.Write(data)

	return err
}

func (c *httpConfig) validateMethod() bool {
	switch c.verb {
	case "GET":

	}
	_, ok := methods[c.verb]

	return ok
}
