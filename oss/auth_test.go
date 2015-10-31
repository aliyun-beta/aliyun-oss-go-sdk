package oss

import (
	"net/http"
	"os"
	"testing"
)

func TestAuth(t *testing.T) {
	req, err := http.NewRequest("PUT", "/oss-example/nelson", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/html")
	req.Header.Set("Date", "Thu, 17 Nov 2005 18:49:58 GMT")
	req.Header.Set("Host", "oss-example.oss-cn-hangzhou.aliyuncs.com")
	req.Header.Set("X-OSS-Meta-Author", "foo@bar.com")
	req.Header.Set("X-OSS-Magic", "abracadabra")
	req.Header.Set("Content-Md5", "ODBGOERFMDMzQTczRUY3NUE3NzA5QzdFNUYzMDQxNEM=")
	auth := authorization{req: req, secret: []byte("OtxrzxIsfpFjA7SwPzILwy8Bw21TLhquhboDYROV")}
	if actual, expected := string(auth.canonicalizedOSSHeaders()),
		"x-oss-magic:abracadabra\nx-oss-meta-author:foo@bar.com\n"; actual != expected {
		t.Fatalf(expectBut, expected, actual)
	}
	if actual, expected := string(auth.data()),
		"PUT\nODBGOERFMDMzQTczRUY3NUE3NzA5QzdFNUYzMDQxNEM=\ntext/html\nThu, 17 Nov 2005 18:49:58 GMT\nx-oss-magic:abracadabra\nx-oss-meta-author:foo@bar.com\n/oss-example/nelson"; actual != expected {
		t.Fatalf(expectBut, expected, actual)
	}
	if actual, expected := auth.value(), "26NBxoKdsyly4EDv6inkoDft/yA="; actual != expected {
		t.Fatalf(expectBut, expected, actual)
	}
}

func TestContentMD5(t *testing.T) {
	f, err := os.Open("testdata/h-content-md5.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	req, err := http.NewRequest("GET", "/", f)
	if err != nil {
		t.Fatal(err)
	}
	setContentMD5(req)
	if expected, actual := "0TMnkhCZtrIjdTtJk6x3+Q==", req.Header.Get("Content-Md5"); actual != expected {
		t.Fatalf(expectBut, expected, actual)
	}
}
