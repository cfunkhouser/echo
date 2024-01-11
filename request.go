// Package echo contains utilities to inspect HTTP requests.
package echo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"

	_ "embed"
)

// Header of an HTTP request.
type Header struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
	Index int    `json:"index"`
	Count int    `json:"repeat_count"`
}

// DumpeHeader from an HTTP request.
func DumpHeader(req *http.Request) (ret []*Header) {
	h := req.Header
	var hns []string
	for n := range h {
		hns = append(hns, n)
	}

	// Sort the headers in the order we want them. Host will always be first.
	sort.Slice(hns, func(i, j int) bool {
		if http.CanonicalHeaderKey(hns[i]) == "Host" {
			return true
		}
		return hns[i] < hns[j]
	})

	for _, n := range hns {
		vs := h.Values(n)
		for i, v := range vs {
			ret = append(ret, &Header{
				Name:  n,
				Value: v,
				Index: i,
				Count: len(vs),
			})
		}
	}
	return
}

// Body of an HTTP payload.
type Body struct {
	// Content of the HTTP body, represented as bytes.
	Content []byte

	// AllegedType of the content, as reported by the client. No validation is
	// done on this value, so be sure you trust the client before relying on it.
	AllegedType string
}

// DumpBody for rendering in an HTML payload.
func DumpBody(req *http.Request) (ret Body) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, req.Body); err != nil {
		ret.Content = []byte(
			fmt.Sprintf("encountered error while dumping body: %v", err))
	}
	defer req.Body.Close()

	ret.AllegedType = req.Header.Get("Content-Type")
	ret.Content = buf.Bytes()
	return
}

// RequestDump is a renderable dump of an inbound HTTP request.
type RequestDump struct {
	Requestor string    `json:"requestor"`
	Method    string    `json:"method"`
	Target    string    `json:"target"`
	Proto     string    `json:"protocol"`
	Headers   []*Header `json:"headers"`
	Body      Body      `json:"-"`
}

// DumpRequest for rendering in an HTML payload.
func DumpRequest(req *http.Request) *RequestDump {
	if req == nil {
		return &RequestDump{}
	}

	return &RequestDump{
		Requestor: req.RemoteAddr,
		Method:    req.Method,
		Target:    req.RequestURI,
		Proto:     req.Proto,
		Headers:   DumpHeader(req),
		Body:      DumpBody(req),
	}
}
