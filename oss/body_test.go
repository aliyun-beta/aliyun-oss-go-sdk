package oss

import (
	"bytes"
	"net/http"
	"os"
	"testing"
)

func TestHTTPBody(t *testing.T) {
	{
		body := "abc"
		req, _ := http.NewRequest("GET", "", nil)
		w := bytes.NewBuffer([]byte(body))
		httpBody(w)(req)
		if int(req.ContentLength) != len(body) {
			t.Fatalf(expectBut, len(body), req.ContentLength)
		}
	}
	{
		body := "abcd"
		req, _ := http.NewRequest("GET", "", nil)
		w := bytes.NewReader([]byte(body))
		httpBody(w)(req)
		if int(req.ContentLength) != len(body) {
			t.Fatalf(expectBut, len(body), req.ContentLength)
		}
	}
	{
		w, err := os.Open("testdata/test")
		if err != nil {
			t.Fatal(err)
		}
		defer w.Close()
		req, _ := http.NewRequest("GET", "", nil)
		httpBody(w)(req)
		if int(req.ContentLength) != 16 {
			t.Fatalf(expectBut, 16, req.ContentLength)
		}
	}
}
