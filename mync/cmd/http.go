package cmd

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type httpConfig struct {
	url      string
	verb     string
	filePath string
	body     string
}

func HandleHttp(w io.Writer, args []string) error {
	hc := httpConfig{}
	fs := flag.NewFlagSet("http", flag.ContinueOnError)

	fs.SetOutput(w)
	fs.StringVar(&hc.verb, "verb", "GET", "HTTP method")
	fs.StringVar(&hc.filePath, "output", "", "File path where save the data's output")
	fs.StringVar(&hc.body, "body", "", "JSON data to be used as payload")
	fs.StringVar(&hc.body, "body-file", "", "File path containing the JSON data to be used as payload")

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

		err = hc.validateMethod(file)
	} else {
		err = hc.validateMethod(w)
	}

	if err != nil {
		return err
	}
	return nil
}

func (hc *httpConfig) validateMethod(w io.Writer) error {
	var data []byte
	var err error
	switch hc.verb {
	case "GET":
		data, err = hc.handleGet()
	case "POST":
		data, err = hc.handlePost()
	default:
		return InvalidHttpMethod
	}

	if err != nil {
		return err
	}
	w.Write(data)

	return nil
}

func (hc *httpConfig) handleGet() ([]byte, error) {
	res, err := http.Get(hc.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func (hc *httpConfig) handlePost() ([]byte, error) {
	if !json.Valid([]byte(hc.body)) {
		return nil, InvalidJsonBody
	} else {
		body := bytes.NewReader([]byte(hc.body))
		res, err := http.Post(hc.url, "application/json", body)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()
		return io.ReadAll(res.Body)
	}

}
