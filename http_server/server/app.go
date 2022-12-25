package server

import (
	"context"
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

type requestIDKey struct{}
type requestCtxValue struct {
	requestID string
}
type logLine struct {
	UserIP string `json:"user_ip"`
	Event  string `json:"event"`
}
type App struct {
	Address string
}

func setupHandlers(sm *http.ServeMux) {
	sm.HandleFunc("/api", apiHandler)
	sm.HandleFunc("/healthz", healthCheckHandler)
	sm.HandleFunc("/decode", decodeHandler)
	sm.HandleFunc("/download", downloadHandler)
	sm.HandleFunc("/job", longRunningProcessHandler)
}

func logMiddleware(req *http.Request) *http.Request {
	l := log.Default()
	v := req.Context().Value(requestIDKey{})

	if m, ok := v.(requestCtxValue); ok {
		msg := struct {
			RequestID string
			Url       string
			Method    string
			BodySize  int64
			Protocol  string
		}{
			m.requestID,
			req.URL.String(),
			req.Method,
			req.ContentLength,
			req.Proto,
		}
		j, err := json.Marshal(&msg)
		if err != nil {
			panic(err)
		}

		l.Printf("%s", j)
	}
	return nil
}

func apiHandler(res http.ResponseWriter, req *http.Request) {
	requestID := "request-123-abc"
	r := addRequestID(req, requestID)

	logMiddleware(r)
	fmt.Fprintln(res, "Hello, world!")
}

func healthCheckHandler(res http.ResponseWriter, req *http.Request) {
	requestID := "request-123-abc"
	r := addRequestID(req, requestID)

	logMiddleware(r)
	fmt.Fprintf(res, "ok")
}

func decodeHandler(res http.ResponseWriter, req *http.Request) {
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

func downloadHandler(res http.ResponseWriter, req *http.Request) {
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

func longRunningProcessHandler(res http.ResponseWriter, req *http.Request) {
	done := make(chan struct{})
	logReader, logWriter := io.Pipe()
	go longRunningProcess(logWriter)
	go progressStreamer(logReader, res, done)

	<-done
}

func addRequestID(r *http.Request, requestID string) *http.Request {
	rq := r.Context()
	ctx := context.WithValue(rq, requestIDKey{}, requestCtxValue{requestID})

	return r.WithContext(ctx)
}

func (app *App) Start() error {
	sm := http.NewServeMux()

	setupHandlers(sm)

	return http.ListenAndServe(app.Address, sm)
}
