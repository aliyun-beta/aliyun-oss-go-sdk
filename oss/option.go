package oss

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type (
	Option  func(*http.Request) error
	HeaderT struct {
		ContentMD5 Option
	}
)

var Header = &HeaderT{
	ContentMD5: setContentMD5,
}

func (*HeaderT) ACL(acl ACL) Option {
	return setHeader("X-Oss-Acl", string(acl))
}
func (*HeaderT) ContentType(value string) Option {
	return setHeader("Content-Type", value)
}
func (*HeaderT) CacheControl(value string) Option {
	return setHeader("Cache-Control", value)
}
func (*HeaderT) ContentDisposition(value string) Option {
	return setHeader("Content-Disposition", value)
}
func (*HeaderT) ContentEncoding(value string) Option {
	return setHeader("Content-Encoding", value)
}
func (*HeaderT) Expires(value string) Option {
	return setHeader("Expires", value)
}
func (*HeaderT) Meta(key, value string) Option {
	return setHeader("X-Oss-Meta-"+key, value)
}
func setHeader(key, value string) Option {
	return func(req *http.Request) error {
		req.Header.Set(key, value)
		return nil
	}
}

func Body(body io.Reader) Option {
	return func(req *http.Request) error {
		rc, ok := body.(io.ReadCloser)
		if !ok && body != nil {
			rc = ioutil.NopCloser(body)
		}
		req.Body = rc
		if body != nil {
			switch v := body.(type) {
			case *bytes.Buffer:
				req.ContentLength = int64(v.Len())
			case *bytes.Reader:
				req.ContentLength = int64(v.Len())
			case *strings.Reader:
				req.ContentLength = int64(v.Len())
			}
		}
		return nil
	}
}
