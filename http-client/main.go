package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type pkgData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type pkgRegisterResult struct {
	ID string `json:"id"`
}

type pkgData2 struct {
	Name     string
	Version  string
	Filename string
	Bytes    io.Reader
}

func registerPackageData(url string, data pkgData) (pkgRegisterResult, error) {
	p := pkgRegisterResult{}
	b, err := json.Marshal(data)
	if err != nil {
		return p, err
	}

	reader := bytes.NewReader(b)
	r, err := http.Post(url, "application/json", reader)

	if err != nil {
		return p, err
	}
	defer r.Body.Close()

	respData, err := io.ReadAll(r.Body)
	if err != nil {
		return p, err
	}

	if r.StatusCode != http.StatusOK {
		return p, errors.New(string(respData))
	}

	err = json.Unmarshal(respData, &p)

	return p, err
}

func fetchPackageData(url string) ([]pkgData, error) {
	var packages []pkgData
	r, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if r.Header.Get("Content-Type") != "application/json" {
		return packages, nil
	}
	data, err := io.ReadAll(r.Body)

	if err != nil {
		return packages, err
	}

	err = json.Unmarshal(data, &packages)

	return packages, err
}

func createFormField(data pkgData2) ([]byte, string, error) {
	var b bytes.Buffer
	var err error
	var fw io.Writer

	mw := multipart.NewWriter(&b)
	fw, err = mw.CreateFormField("name")

	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Name)

	fw, err = mw.CreateFormField("version")
	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Version)

	fw, err = mw.CreateFormFile("filedata", data.Filename)
	if err != nil {
		return nil, "", err
	}

	_, err = io.Copy(fw, data.Bytes)
	err = mw.Close()
	if err != nil {
		return nil, "", err
	}

	contentType := mw.FormDataContentType()

	return b.Bytes(), contentType, nil
}

func main() {
	client := http.Client{
		Timeout: 100 * time.Millisecond,
	}
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "Must specify a HTTP URL to get data from")
		os.Exit(1)
	}

	body, err := fetchRemoteResource(&client, os.Args[1])

	if err != nil {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s\n", body)
}

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	return &http.Client{Timeout: d}
}

func fetchRemoteResource(client *http.Client, url string) ([]byte, error) {
	res, err := client.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return io.ReadAll(res.Body)
}
