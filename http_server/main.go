package main

import (
	"http_server/server"
	"log"
	"os"
)

func main() {
	listenAddr := os.Getenv("Listen_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	app := server.App{
		Address: listenAddr,
		Logger:  log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)}

	log.Fatal(app.Start())
}
