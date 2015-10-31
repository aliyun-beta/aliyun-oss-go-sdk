package oss

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
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
		if f, ok := body.(*os.File); ok {
			fInfo, err := os.Stat(f.Name())
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", mime.TypeByExtension(path.Ext(f.Name())))
			req.ContentLength = fInfo.Size()
		}
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
