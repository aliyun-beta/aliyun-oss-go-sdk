package oss

import (
	"net/http"
	"testing"
)

func TestSign(t *testing.T) {
	req, err := http.NewRequest("PUT", "/nelson", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/html")
	req.Header.Set("Date", "Thu, 17 Nov 2005 18:49:58 GMT")
	req.Header.Set("Host", "oss-example.oss-cn-hangzhou.aliyuncs.com")
	req.Header.Set("X-OSS-Meta-Author", "foo@bar.com")
	req.Header.Set("X-OSS-Magic", "abracadabra")
	req.Header.Set("Content-Md5", "ODBGOERFMDMzQTczRUY3NUE3NzA5QzdFNUYzMDQxNEM=")
	sig := signature{req}
	if actual, expected := string(sig.canonicalizedOSSHeaders()), "x-oss-magic:abracadabra\nx-oss-meta-author:foo@bar.com"; actual != expected {
		t.Fatalf(expectBut, expected, actual)
	}
}
