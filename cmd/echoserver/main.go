// Program echoserver is an HTTP server which accepts any request and responds
// with with 200 OK. The response payload is a web page which displays
// information about the request and the server environment.
package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"idontfixcomputers.com/echo"
)

var address = flag.String("address", ":8080",
	"Bind address in host:port format.")

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	flag.Parse()
	if *address == "" {
		logger.Error("Can't bind an empty address")
		os.Exit(1)
	}

	http.Handle("/", echo.New(logger))
	logger.Info("Server starting", "address", *address)

	err := http.ListenAndServe(*address, nil)
	logger.Info("Server has exited", "error", err)
}
