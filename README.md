# echo - HTTP Server Environment Debugger

`echo` accepts any request and responds with with `200 OK`. The response payload
is a web page which displays information about the request and the server
environment.

All requests to the server are logged as JSON.

## The Server

You can run the server directly from this repository:

```console
$ go run ./cmd/echoserver
INFO[0000] Server starting                               address=":8080"
```

You can install it directly, and run it:

```console
$ go install idontfixcomputers.com/echo/cmd/echoserver@latest
go: downloading idontfixcomputers.com/echo v0.1.1
$ which echo
/home/whatever/bin/echo
$ echo --help
Usage of echo:
  -address string
        Bind address in host:port format. (default ":8080")
$ echo --address 127.0.0.1:9000
{"time":"2023-12-20T19:41:36.894162295Z","level":"INFO","msg":"Server starting","address":":9000"}
```

Or, you can run it via Docker with `docker run --rm -p 8080:8080
cfunkhouser/echoserver:latest`.
