package cmd

import (
	"log"
	"net/http"
	"time"
)

type LoggingClient struct {
	log *log.Logger
}

func (c LoggingClient) RoundTrip(r *http.Request) (*http.Response, error) {
	ds := time.Now()
	c.log.Printf(
		"Sending a %s request to %s over %s\n",
		r.Method, r.URL, r.Proto,
	)

	res, err := http.DefaultTransport.RoundTrip(r)
	de := time.Now()
	te := de.Sub(ds)

	c.log.Printf("Request took %f seconds to complete", te.Seconds())
	c.log.Printf("Got back a response over %s\n", res.Proto)

	return res, err
}
