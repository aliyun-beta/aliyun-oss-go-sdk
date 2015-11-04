package oss

import (
	"bytes"
	"errors"
	"mime/multipart"
	"testing"
)

type postOptionTestCase struct {
	option PostOption
	key    string
	value  string
}

var postOptionTestcases = []postOptionTestCase{
	{
		option: PostContentType("plain/text"),
		key:    "Content-Type",
		value:  "plain/text",
	},
	{
		option: PostCacheControl("no-cache"),
		key:    "Cache-Control",
		value:  "no-cache",
	},
	{
		option: PostContentDisposition("Attachment; filename=example.txt"),
		key:    "Content-Disposition",
		value:  "Attachment; filename=example.txt",
	},
	{
		option: PostContentEncoding("gzip"),
		key:    "Content-Encoding",
		value:  "gzip",
	},
	{
		option: PostExpires("Thu, 01 Dec 1994 16:00:00 GMT"),
		key:    "Expires",
		value:  "Thu, 01 Dec 1994 16:00:00 GMT",
	},
	{
		option: PostSuccessActionRedirect("http://example.com"),
		key:    "success_action_redirect",
		value:  "http://example.com",
	},
	{
		option: PostServerSideEncryption("AES256"),
		key:    "x-oss-server-side-encryption",
		value:  "AES256",
	},
	{
		option: PostObjectACL(PublicReadWriteACL),
		key:    "x-oss-object-acl",
		value:  "public-read-write",
	},
}

func TestPostOptions(t *testing.T) {
	boundary := "abc"
	for _, testcase := range postOptionTestcases {
		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		w.SetBoundary(boundary)
		if err := testcase.option(w); err != nil {
			t.Fatalf(testcaseErr, testcase.key, err)
		}
		w.Close()
		form, err := multipart.NewReader(buf, boundary).ReadForm(8192)
		if err != nil {
			t.Fatalf(testcaseErr, testcase.key, err)
		}
		if vs := form.Value[testcase.key]; len(vs) != 1 || vs[0] != testcase.value {
			t.Fatalf(testcaseExpectBut, testcase.key, testcase.value, vs)
		}
	}
}

func TestFailedPostOption(t *testing.T) {
	injectedFailure := errors.New("injected failure")
	failedOption := func(*multipart.Writer) error {
		return injectedFailure
	}
	api := New(testEndpoint, testID, testSecret)
	if _, err := api.PostObject(testBucketName, testObjectName, "filename", "", failedOption); err != injectedFailure {
		t.Fatalf(expectBut, injectedFailure, err)
	}
}

func TestPostFileNotExists(t *testing.T) {
	filename := "fjssfsefwel324234ufj32"
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	defer w.Close()
	if err := postFile(filename)(w); err == nil {
		t.Fatal("expect error but got nil")
	}
}

func TestSafeCopy(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := safeCopy(nil, buf); err == nil {
		t.Fatal("expect error but got nil")
	}
	if err := safeCopy(buf, nil); err == nil {
		t.Fatal("expect error but got nil")
	}
}
