package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type pkgRegisterResult struct {
	Id string `json:"id"`
}
type pkgData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func main() {

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
		return nil, err
	}
	err = json.Unmarshal(data, &packages)
	if err != nil {
		return nil, err
	}
	return packages, nil
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
