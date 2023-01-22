package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func startTestPackageServer() *httptest.Server {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" {
					pkgData := `[
{"name": "package1", "version": "1.1"},
{"name": "package2", "version": "1.0"}
]`
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, pkgData)
				}
				if r.Method == "POST" {
					d := pkgRegisterResult{}
					err := r.ParseMultipartForm(5000)
					if err != nil {
						http.Error(
							w, err.Error(), http.StatusBadRequest,
						)
						return
					}
					mForm := r.MultipartForm
					f := mForm.File["filedata"][0]
					d.Id = fmt.Sprintf(
						"%s-%s", mForm.Value["name"][0], mForm.Value["version"][0],
					)
					d.Filename = f.Filename
					d.Size = f.Size
					jsonData, err := json.Marshal(d)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, string(jsonData))
				}
			}))

	return ts
}

func packageRegHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Incoming package data
		p := pkgData{}

		// Package registration response
		d := pkgRegisterResult{}
		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(data, &p)
		if err != nil || len(p.Name) == 0 || len(p.Version) == 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		d.Id = p.Name + "-" + p.Version
		jsonData, err := json.Marshal(d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonData))
	} else {
		http.Error(w, "Invalid HTTP method specified", http.StatusMethodNotAllowed)
		return
	}
}

func TestFetchPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()

	packages, err := fetchPackageData(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if len(packages) != 2 {
		t.Fatalf("Expected 2 packages, Got back: %d", len(packages))
	}
}

func TestRegisterPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()
	p := pkgData{
		Name:     "mypackage",
		Version:  "0.1",
		Filename: "mypackage-0.1.tar.gz",
		Bytes:    strings.NewReader("data"),
	}
	client := createHTTPClientWithTimeout(time.Second * 5)
	pResult, err := registerPackageData(client, ts.URL, p)
	if err != nil {
		t.Fatal(err)
	}

	if pResult.Id != fmt.Sprintf("%s-%s", p.Name, p.Version) {
		t.Errorf(
			"Expected package ID to be %s-%s, Got: %s", p.Name, p.Version, pResult.Id,
		)
	}
	if pResult.Filename != p.Filename {
		t.Errorf(
			"Expected package filename to be %s, Got: %s", p.Filename, pResult.Filename,
		)
		if pResult.Size != 4 {
			t.Errorf("Expected package size to be 4, Got: %d", pResult.Size)
		}
	}
}
