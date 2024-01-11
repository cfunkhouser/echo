package echo

import (
	"bytes"
	"io"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func discardLogger(tb testing.TB) *slog.Logger {
	tb.Helper()
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func bodyWriter(tb testing.TB, body string) io.ReadCloser {
	var buf bytes.Buffer
	if _, err := (&buf).WriteString(body); err != nil {
		tb.Fatalf("error setting up test: %v", err)
	}
	return io.NopCloser(&buf)
}

func TestEchoHandlerServeHTTP(t *testing.T) {
	r := httptest.NewRequest("GET", "https://example.com:9443/api", bodyWriter(t, "Hello there!"))
	r.Header.Add("X-Some-Thing", "testing")

	w := httptest.NewRecorder()

	(&EchoHandler{
		Log: discardLogger(t),
	}).ServeHTTP(w, r)

	resp := w.Result()

	if got := resp.StatusCode; got != 200 {
		t.Errorf("ServeHTTP(): StatusCode mismatch: got: %v, want: 200", got)
	}

	wantCTH := "text/plain; charset=utf-8"
	if got := resp.Header.Get("Content-Type"); got != wantCTH {
		t.Errorf("ServeHTTP(): Content-Type Header mismatch: got: %q, want: %q", got, wantCTH)
	}

	got, _ := io.ReadAll(resp.Body)
	want := "GET https://example.com:9443/api HTTP/1.1\r\nX-Some-Thing: testing\r\n\r\nHello there!"
	if diff := cmp.Diff(want, string(got)); diff != "" {
		t.Errorf("ServeHTTP(): response body mismatch (-want +got):\n%v", diff)
	}
}
