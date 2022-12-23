// chap5/basic-http-server/server.go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

/*
Notes

- How handler functions?
When the server receives a request for a path, for example /api/, the struct ServerMux has a field that maps the path with a correspodent handler function
This handler functions has a interface, like func(res http.Writer, req *http.Request) error, that not return anything, instead it respond to the client via the http.Writter interface.
For each comming client request, http package creates a new goroutine to handler the request. Doing that allows each behavior request don't interfer to any others.
- How processing request?
The “*http.Request“ has several field with information about the request and methods to process it. The most important are Method (string), URL (*url.URL), Proto (String protocol), Header (map[string][]string), Host, Body (That is io.ReadCloser, thus any function that accepts a io.Reader can be used to handle this body, as example io.ReadAll), Form PostForm, MultipartForm
For each incoming request processed by a handler function, a *context* is associated for it. The context can be obtained by calling “req.Context()“. The life cycle of this context is the same as that of the request.
- How read and write streaming data, that is, how send data while it becomes available?
Normally, a response of a http request is send quickly, about a few seconds. However, there are long running jobs that take a time to process the request. In that case, the connection between the client and server must be active and the server must be send a chunck of data while it becomes available.
For that, uses the “io.Pipe“ that returns a “io.PipeWriter, io.PipeReader“. The producer will be send that to “io.PipeWriter“ and the consumer gonna read the “io.PipeReader“ to get the data and send back to the client.
To limit the length of chunk to be send, we can use the http.Flush if the “io.ResponseWriter“ implement the flush interface.
*/

func setupHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthCheckHandler)
	mux.HandleFunc("/api", apiHandler)
	mux.HandleFunc("/decode", decodeHandler)
	mux.HandleFunc("/job", longRunningProcessHandler)
}

func main() {
	listenAddr := os.Getenv("Listen_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	setupHandlers(mux)
	log.Fatal(http.ListenAndServe(listenAddr, mux))
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

type logLine struct {
	UserIP string `json:"user_ip"`
	Event  string `json:"event"`
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	js := json.NewDecoder(r.Body)
	js.DisallowUnknownFields()
	var e *json.UnmarshalTypeError

	for {
		var l logLine

		err := js.Decode(&l)
		if err == io.EOF {
			break
		}
		if errors.As(err, &e) {
			log.Println(err)
			continue
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(l.UserIP, l.Event)
	}
	fmt.Fprintf(w, "OK")
}

func printLog(req *http.Request) {
	l := log.Default()
	ctx := req.Context()
	v := ctx.Value(requestContextKey{})

	if m, ok := v.(requestContextValue); ok {
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
}

type requestContextKey struct{}
type requestContextValue struct {
	requestID string
}

// helper function to store a request ID in a request's context
func addRequestID(r *http.Request, requestID string) *http.Request {
	c := requestContextValue{
		requestID,
	}
	currentCtx := r.Context()
	newCtx := context.WithValue(currentCtx, requestContextKey{}, c)

	return r.WithContext(newCtx)

}

func apiHandler(res http.ResponseWriter, req *http.Request) {
	requestID := "request-123-abc"
	r := addRequestID(req, requestID)

	printLog(r)
	fmt.Fprintf(res, "Hello, world!")
}

func healthCheckHandler(res http.ResponseWriter, req *http.Request) {
	requestID := "request-123-abc"
	r := addRequestID(req, requestID)

	printLog(r)
	fmt.Fprintf(res, "ok")
}
