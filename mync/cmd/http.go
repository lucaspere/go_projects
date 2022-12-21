package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type httpConfig struct {
	url            string
	verb           string
	outputFilePath string
	body           string
	filePath       string
	formData       map[string]string
	disableRed     bool
	headers        map[string]string
	basicAuth      string
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
	fs.StringVar(&hc.basicAuth, "basicauth", "", "Basic authentication for the request")
	fs.BoolVar(&hc.disableRed, "disable-redirect", false, "If it is for the client not to follow the redirect url")

	fd := fs.String("form-data", "", "Form-Data key-value pair")
	h := fs.String("header", "", "HTTP request header")

	if len(*fd) > 0 {
		form := strings.Split(*fd, "=")

		hc.formData[form[0]] = form[1]
	}
	if len(*h) > 0 {
		form := strings.Split(*h, "=")

		hc.headers[form[0]] = form[1]
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
	myTransport := LoggingClient{}
	l := log.New(os.Stdout, "", log.LstdFlags)
	myTransport.log = l
	client := http.Client{
		CheckRedirect: hc.redirectPolicy,
		Transport:     &myTransport,
	}

	req, err := hc.createHTTPGetRequest()
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func (hc *httpConfig) handlePost() ([]byte, error) {
	var body io.Reader
	var ct = "application/json"
	client := http.Client{
		CheckRedirect: hc.redirectPolicy,
	}
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

	res, err := client.Post(hc.url, ct, body)
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

		fmt.Fprintln(fw, value)
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

func (hc *httpConfig) redirectPolicy(req *http.Request, via []*http.Request) error {
	if hc.disableRed {
		return errors.New("disabled redirect")
	}
	return nil
}

func (hc *httpConfig) createHTTPGetRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", hc.url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range hc.headers {
		req.Header.Add(k, v)
	}

	if len(hc.basicAuth) > 0 {
		err = hc.setAuth(req)
	}

	return req, err
}

func (hc *httpConfig) setAuth(req *http.Request) error {
	if r, e := regexp.Match(`\w+:\w+`, []byte("")); e != nil || !r {
		return errors.New("invalid basicauth value")
	}
	auth := strings.Split(hc.basicAuth, ":")
	req.SetBasicAuth(auth[0], auth[1])

	return nil
}
