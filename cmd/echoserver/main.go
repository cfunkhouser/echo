// Program echoserver is an HTTP server which accepts any request and responds
// with with 200 OK. The response payload is a web page which displays
// information about the request and the server environment.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"idontfixcomputers.com/echo"
)

func envOr(env, def string) string {
	port := def
	if v, ok := os.LookupEnv(env); ok {
		port = v
	}
	return fmt.Sprintf(":%s", port)
}

var address = flag.String("address", envOr("PORT", "8080"),
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
