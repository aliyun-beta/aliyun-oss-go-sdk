package oss

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type (
	Option func(*http.Request) error
)

func ACL(acl ACLType) Option {
	return setHeader("X-Oss-Acl", string(acl))
}
func ContentType(value string) Option {
	return setHeader("Content-Type", value)
}
func CacheControl(value string) Option {
	return setHeader("Cache-Control", value)
}
func ContentDisposition(value string) Option {
	return setHeader("Content-Disposition", value)
}
func ContentEncoding(value string) Option {
	return setHeader("Content-Encoding", value)
}
func Expires(value string) Option {
	return setHeader("Expires", value)
}
func ResponseContentType(value string) Option {
	return setHeader("Response-Content-Type", value)
}
func ResponseContentLanguage(value string) Option {
	return setHeader("Response-Content-Language", value)
}
func ResponseCacheControl(value string) Option {
	return setHeader("Response-Cache-Control", value)
}
func ResponseContentDisposition(value string) Option {
	return setHeader("Response-Content-Disposition", value)
}
func ResponseContentEncoding(value string) Option {
	return setHeader("Response-Content-Encoding", value)
}
func ResponseExpires(value string) Option {
	return setHeader("Response-Expires", value)
}
func Meta(key, value string) Option {
	return setHeader("X-Oss-Meta-"+key, value)
}
func Range(value string) Option {
	return setHeader("Range", value)
}
func IfModifiedSince(value string) Option {
	return setHeader("If-Modified-Since", value)
}
func IfUnmodifiedSince(value string) Option {
	return setHeader("If-Unmodified-Since", value)
}
func IfMatch(value string) Option {
	return setHeader("If-Match", value)
}
func IfNoneMatch(value string) Option {
	return setHeader("If-None-Match", value)
}
func setHeader(key, value string) Option {
	return func(req *http.Request) error {
		req.Header.Set(key, value)
		return nil
	}
}

func httpBody(body io.Reader) Option {
	return func(req *http.Request) error {
		if f, ok := body.(*os.File); ok {
			fInfo, err := os.Stat(f.Name())
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", typeByExtension(f.Name()))
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
func typeByExtension(file string) string {
	typ := mime.TypeByExtension(path.Ext(file))
	if typ == "" {
		typ = "application/octet-stream"
	}
	return typ
}

func Delimiter(value string) Option {
	return addParam("delimiter", value)
}
func Marker(value string) Option {
	return addParam("marker", value)
}
func MaxKeys(value int) Option {
	return addParam("maxkeys", strconv.Itoa(value))
}
func Prefix(value string) Option {
	return addParam("prefix", value)
}
func EncodingType(value string) Option {
	return addParam("encoding-type", value)
}
func addParam(key, value string) Option {
	return func(req *http.Request) error {
		q := req.URL.Query()
		q.Add(key, value)
		req.URL.RawQuery = q.Encode()
		return nil
	}
}

type CreateBucketConfiguration struct {
	LocationConstraint string
}

func BucketLocation(value string) Option {
	return func(req *http.Request) error {
		var w bytes.Buffer
		w.WriteString(xml.Header)
		if err := xml.NewEncoder(&w).Encode(CreateBucketConfiguration{LocationConstraint: value}); err != nil {
			return err
		}
		return httpBody(bytes.NewReader(w.Bytes()))(req)
	}
}

type Delete struct {
	Quiet  bool
	Object []Object
}
type Object struct {
	Key string
}

func deleteObjects(objects []string, quiet bool) Option {
	return func(req *http.Request) error {
		var w bytes.Buffer
		w.WriteString(xml.Header)
		del := Delete{Quiet: quiet}
		for _, key := range objects {
			del.Object = append(del.Object, Object{Key: key})
		}
		if err := xml.NewEncoder(&w).Encode(del); err != nil {
			return err
		}
		return httpBody(bytes.NewReader(w.Bytes()))(req)
	}
}
