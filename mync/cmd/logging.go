package cmd

import (
	"log"
	"net/http"
)

type LoggingClient struct {
	log *log.Logger
}

func (c LoggingClient) RoundTrip(r *http.Request) (*http.Response, error) {
	c.log.Printf(
		"Sending a %s request to %s over %s\n",
		r.Method, r.URL, r.Proto,
	)

	res, err := http.DefaultTransport.RoundTrip(r)
	c.log.Printf("Got back a response over %s\n", res.Proto)

	return res, err
}
