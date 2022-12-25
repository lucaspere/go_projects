package server

import (
	"log"
	"net/http"

	c "http_server/server/controllers"
	md "http_server/server/middlewares"
)

type App struct {
	Address string
	Logger  *log.Logger
}

func (a *App) setupHandlers(sm *http.ServeMux) {
	sm.HandleFunc("/api", c.ApiHandler)
	sm.HandleFunc("/healthz", c.HealthCheckHandler)
	sm.HandleFunc("/decode", c.DecodeHandler)
	sm.HandleFunc("/download", c.DownloadHandler)
	sm.HandleFunc("/job", c.LongRunningProcessHandler)
}

func (app *App) Start() error {
	sm := http.NewServeMux()

	app.setupHandlers(sm)

	m := md.AddRequestID(md.LoggingMiddleware(md.PanicMiddleware(sm)))

	return http.ListenAndServe(app.Address, m)
}
