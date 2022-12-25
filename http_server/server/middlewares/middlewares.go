package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
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

func AddRequestID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := rand.NewSource(time.Now().Unix())
		rand := rand.New(s)
		requestID := "request-id-" + fmt.Sprint(rand.Int())
		rc := r.Context()
		ctx := context.WithValue(rc, requestIDKey{}, requestCtxValue{requestID})

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			h.ServeHTTP(w, r)
			l := log.Default()
			v := r.Context().Value(requestIDKey{})

			if m, ok := v.(requestCtxValue); ok {
				msg := struct {
					RequestID string
					Path      string
					Method    string
					BodySize  int64
					Protocol  string
					Duration  float64
				}{
					m.requestID,
					r.URL.String(),
					r.Method,
					r.ContentLength,
					r.Proto,
					time.Since(startTime).Seconds(),
				}
				j, _ := json.Marshal(&msg)

				l.Printf("%s", j)
			}
		},
	)
}

func PanicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rValue := recover(); rValue != nil {
				log.Println("panic detected", rValue)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Unexpected server error")
			}
		}()

		h.ServeHTTP(w, r)
	})
}
