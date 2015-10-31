package oss

import (
	"net/http"
	"testing"
)

type optionTestCase struct {
	option Option
	key    string
	value  string
}

var headerTestcases = []optionTestCase{
	{
		option: Meta("User", "baymax"),
		key:    "X-Oss-Meta-User",
		value:  "baymax",
	},
	{
		option: ACL(PrivateACL),
		key:    "X-Oss-Acl",
		value:  "private",
	},
	{
		option: ContentType("plain/text"),
		key:    "Content-Type",
		value:  "plain/text",
	},
	{
		option: CacheControl("no-cache"),
		key:    "Cache-Control",
		value:  "no-cache",
	},
	{
		option: ContentDisposition("Attachment; filename=example.txt"),
		key:    "Content-Disposition",
		value:  "Attachment; filename=example.txt",
	},
	{
		option: ContentEncoding("gzip"),
		key:    "Content-Encoding",
		value:  "gzip",
	},
	{
		option: Expires("Thu, 01 Dec 1994 16:00:00 GMT"),
		key:    "Expires",
		value:  "Thu, 01 Dec 1994 16:00:00 GMT",
	},
}

func TestHeaderOptions(t *testing.T) {
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

var paramTestCases = []optionTestCase{
	{
		option: Delimiter("/"),
		key:    "delimiter",
		value:  "/",
	},
	{
		option: Marker("abc"),
		key:    "marker",
		value:  "abc",
	},
	{
		option: MaxKeys(150),
		key:    "maxkeys",
		value:  "150",
	},
	{
		option: Prefix("fun"),
		key:    "prefix",
		value:  "fun",
	},
	{
		option: EncodingType("ascii"),
		key:    "encoding-type",
		value:  "ascii",
	},
}

func TestParamOptions(t *testing.T) {
	for i, testcase := range paramTestCases {
		req, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatal("testcase", i, err)
		}
		if err := testcase.option(req); err != nil {
			t.Fatal("testcase", i, err)
		}
		if expected, actual := testcase.value, req.URL.Query().Get(testcase.key); actual != expected {
			t.Fatalf(testcaseExpectBut, i, expected, actual)
		}
	}
}
