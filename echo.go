package echo

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"os"

	_ "embed"
)

// EchoHandler renders HTML payloads containing the in-bound request along with
// some details about the running environment of the server.
type EchoHandler struct {
	Log *slog.Logger
}

var (
	//go:embed templates/echo.tmpl
	echoTemplateSource string
	echoTemplate       = mustParse("echo.tmpl", echoTemplateSource)
)

func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	d := DumpRequest(r)
	if d == nil {
		http.Error(w, "Failed to dump request!", http.StatusBadRequest)
		h.Log.ErrorContext(ctx, "failed to dump request")
		return
	}
	h.Log.InfoContext(ctx, "incoming request", slog.Any("request", d))

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	if err := echoTemplate.Execute(&buf, render(d, nil)); err != nil {
		http.Error(w, "Failed to render dump!", http.StatusInternalServerError)
		h.Log.ErrorContext(ctx, "failed to render dump", slog.Any("error", err))
		return
	}

	// It's too late to do anything graceful with any error as far as the client
	// is concerned, so just log it.
	if _, err := io.Copy(w, &buf); err != nil {
		h.Log.ErrorContext(ctx, "failed write to client", slog.Any("error", err))
	}
}

var defaultEchoHandler = &EchoHandler{
	Log: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
}

// Echo the request back to the client as plain text.
func Echo(w http.ResponseWriter, r *http.Request) {
	defaultEchoHandler.ServeHTTP(w, r)
}

// New echo handler for use in HTTP servers.
func New(logger *slog.Logger) *EchoHandler {
	return &EchoHandler{
		Log: logger,
	}
}
