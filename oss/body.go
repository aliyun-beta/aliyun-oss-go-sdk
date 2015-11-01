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
	"strings"
)

func httpBody(body io.Reader) Option {
	return func(req *http.Request) error {
		fileName := ""
		if f, ok := body.(*os.File); ok {
			fInfo, err := os.Stat(f.Name())
			if err != nil {
				return err
			}
			fileName = f.Name()
			req.ContentLength = fInfo.Size()
		}
		req.Header.Set("Content-Type", typeByExtension(fileName))
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

type CompleteMultipartUpload struct {
	Part []Part
}

func completeMultipartUpload(list *CompleteMultipartUpload) Option {
	return func(req *http.Request) error {
		var w bytes.Buffer
		w.WriteString(xml.Header)
		if err := xml.NewEncoder(&w).Encode(list); err != nil {
			return err
		}
		return httpBody(bytes.NewReader(w.Bytes()))(req)
	}
}
