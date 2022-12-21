package cmd

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type httpConfig struct {
	url            string
	verb           string
	outputFilePath string
	body           string
	filePath       string
	formData       map[string]string
}

func HandleHttp(w io.Writer, args []string) error {
	hc := httpConfig{}
	fs := flag.NewFlagSet("http", flag.ContinueOnError)

	fs.SetOutput(w)
	fs.StringVar(&hc.verb, "verb", "GET", "HTTP method")
	fs.StringVar(&hc.outputFilePath, "output", "", "File path where save the data's output")
	fs.StringVar(&hc.body, "body", "", "JSON data to be used as payload")
	fs.StringVar(&hc.filePath, "body-file", "", "File path containing the JSON data to be used as payload")
	fs.StringVar(&hc.filePath, "upload", "", "File path of the upload file")
	str := fs.String("form-data", "", "Form-Data key-value pair")

	if len(*str) > 0 {
		form := strings.Split(*str, "=")

		hc.formData[form[0]] = form[1]
	}

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

	if len(hc.outputFilePath) > 0 {
		var file *os.File

		path := filepath.Join(hc.outputFilePath)
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
	var body io.Reader
	var ct = "application/json"
	if len(hc.filePath) > 0 {
		hc.filePath = filepath.Join(hc.filePath)
		file, err := os.Open(hc.filePath)
		if err != nil {
			return nil, err
		}

		j, err := io.ReadAll(file)
		if err != nil {
			return j, err
		}

		if len(hc.formData) > 0 {
			j, ct, err = hc.handleMultipart(j)
			if err != nil {
				return j, err
			}
		}
		if !json.Valid(j) {
			return nil, InvalidJsonBody
		}

		body = bytes.NewBuffer(j)

	} else if len(hc.body) > 0 {
		if !json.Valid([]byte(hc.body)) {
			return nil, InvalidJsonBody
		}

		body = bytes.NewReader([]byte(hc.body))
	}

	res, err := http.Post(hc.url, ct, body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return io.ReadAll(res.Body)
}

func (hc *httpConfig) handleMultipart(b []byte) ([]byte, string, error) {
	var buf bytes.Buffer
	var err error
	var fw io.Writer

	mw := multipart.NewWriter(&buf)

	for key, value := range hc.formData {
		fw, err = mw.CreateFormField(key)
		if err != nil {
			return nil, "", err
		}

		fmt.Fprintf(fw, value)
	}

	fw, err = mw.CreateFormFile("filedata", hc.filePath)
	if err != nil {
		return nil, "", err
	}
	data := bytes.NewReader(b)
	_, err = io.Copy(fw, data)
	err = mw.Close()
	if err != nil {
		return nil, "", err
	}

	contentType := mw.FormDataContentType()
	return buf.Bytes(), contentType, nil
}
