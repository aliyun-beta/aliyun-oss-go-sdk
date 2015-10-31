package oss

import (
	"net/http"
	"testing"
)

type headerTestcase struct {
	option Option
	key    string
	value  string
}

var headerTestcases = []headerTestcase{
	{
		option: Header.Meta("User", "baymax"),
		key:    "X-Oss-Meta-User",
		value:  "baymax",
	},
	{
		option: Header.ACL(PrivateACL),
		key:    "X-Oss-Acl",
		value:  "private",
	},
	{
		option: Header.ContentType("plain/text"),
		key:    "Content-Type",
		value:  "plain/text",
	},
	{
		option: Header.CacheControl("no-cache"),
		key:    "Cache-Control",
		value:  "no-cache",
	},
	{
		option: Header.ContentDisposition("Attachment; filename=example.txt"),
		key:    "Content-Disposition",
		value:  "Attachment; filename=example.txt",
	},
	{
		option: Header.ContentEncoding("gzip"),
		key:    "Content-Encoding",
		value:  "gzip",
	},
	{
		option: Header.Expires("Thu, 01 Dec 1994 16:00:00 GMT"),
		key:    "Expires",
		value:  "Thu, 01 Dec 1994 16:00:00 GMT",
	},
}

func TestOptions(t *testing.T) {
	for i, testcase := range headerTestcases {
		req, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatal("testcase", i, err)
		}
		if err := testcase.option(req); err != nil {
			t.Fatal("testcase", i, err)
		}
		if expected, actual := testcase.value, req.Header.Get(testcase.key); actual != expected {
			t.Fatalf(testcaseExpectBut, i, expected, actual)
		}
	}
}
