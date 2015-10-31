package oss

import (
	"io/ioutil"
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
	{
		option: ResponseContentType("plain/text"),
		key:    "Response-Content-Type",
		value:  "plain/text",
	},
	{
		option: ResponseContentLanguage("en"),
		key:    "Response-Content-Language",
		value:  "en",
	},
	{
		option: ResponseCacheControl("no-cache"),
		key:    "Response-Cache-Control",
		value:  "no-cache",
	},
	{
		option: ResponseContentDisposition("Attachment; filename=example.txt"),
		key:    "Response-Content-Disposition",
		value:  "Attachment; filename=example.txt",
	},
	{
		option: ResponseContentEncoding("gzip"),
		key:    "Response-Content-Encoding",
		value:  "gzip",
	},
	{
		option: ResponseExpires("Thu, 01 Dec 1994 16:00:00 GMT"),
		key:    "Response-Expires",
		value:  "Thu, 01 Dec 1994 16:00:00 GMT",
	},
	{
		option: Range("bytes=0-9"),
		key:    "Range",
		value:  "bytes=0-9",
	},
	{
		option: IfModifiedSince("Fri, 24 Feb 2012 06:38:30 GMT"),
		key:    "If-Modified-Since",
		value:  "Fri, 24 Feb 2012 06:38:30 GMT",
	},
	{
		option: IfUnmodifiedSince("Fri, 24 Feb 2012 06:38:30 GMT"),
		key:    "If-Unmodified-Since",
		value:  "Fri, 24 Feb 2012 06:38:30 GMT",
	},
	{
		option: IfMatch("xyzzy"),
		key:    "If-Match",
		value:  "xyzzy",
	},
	{
		option: IfNoneMatch("xyzzy"),
		key:    "If-None-Match",
		value:  "xyzzy",
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

func TestBucketLocationContraint(t *testing.T) {
	req, err := http.NewRequest("PUT", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := BucketLocation("oss-cn-beijing")(req); err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	if expected, actual := `<?xml version="1.0" encoding="UTF-8"?>
<CreateBucketConfiguration><LocationConstraint>oss-cn-beijing</LocationConstraint></CreateBucketConfiguration>`, string(body); actual != expected {
		t.Fatalf(expectBut, expected, actual)
	}
}
