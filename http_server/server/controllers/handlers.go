package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type logLine struct {
	UserIP string `json:"user_ip"`
	Event  string `json:"event"`
}

func ApiHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Hello, world!")
}

func HealthCheckHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(res, "ok")
}

func DecodeHandler(res http.ResponseWriter, req *http.Request) {
	var ue *json.UnmarshalTypeError

	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()

	for {
		var ll logLine

		err := d.Decode(&ll)
		if err != io.EOF {
			break
		}

		if errors.As(err, &ue) {
			log.Println(err)
			continue
		}

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(ll.UserIP, ll.Event)
	}

	fmt.Fprintf(res, "OK")
}

func DownloadHandler(res http.ResponseWriter, req *http.Request) {
	if q := req.URL.Query().Get("fileName"); len(q) > 0 {
		dir, _ := os.Getwd()
		p := filepath.Join(dir, "/files", q)
		f, err := os.Open(p)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		}
		defer f.Close()

		res.Header().Set("Content-Type", "text/plain")
		res.Header().Set("X-Content-Type-Options", "nosniff")
		io.Copy(res, f)

	}
}

func longRunningProcess(logWriter *io.PipeWriter) {
	for i := 0; i < 20; i++ {
		fmt.Fprintf(
			logWriter,
			`{"id": %d, "user_ip": "172.121.19.21", "event": "click_on_add_cart" }`, i,
		)
		fmt.Fprintln(logWriter)
		time.Sleep(1 * time.Second)
	}
	logWriter.Close()
}

func progressStreamer(logReader *io.PipeReader, w http.ResponseWriter, done chan struct{}) {
	buf := make([]byte, 500)

	f, flushSupported := w.(http.Flusher)

	defer logReader.Close()
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	for {
		n, err := logReader.Read(buf)
		if err == io.EOF {
			break
		}
		w.Write(buf[:n])
		if flushSupported {
			f.Flush()
		}
	}
	done <- struct{}{}
}

func LongRunningProcessHandler(res http.ResponseWriter, req *http.Request) {
	done := make(chan struct{})
	logReader, logWriter := io.Pipe()
	go longRunningProcess(logWriter)
	go progressStreamer(logReader, res, done)

	<-done
}
