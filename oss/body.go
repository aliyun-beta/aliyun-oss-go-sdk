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
	"time"
)

type (
	CreateBucketConfiguration struct {
		LocationConstraint string
	}

	Delete struct {
		Quiet  bool
		Object []Object
	}
	Object struct {
		Key string
	}

	CompleteMultipartUpload struct {
		Part []Part
	}

	CORSConfiguration struct {
		CORSRule []CORSRule
	}
	CORSRule struct {
		AllowedOrigin []string
		AllowedMethod []string
		AllowedHeader []string
		ExposeHeader  []string
		MaxAgeSeconds int `xml:"MaxAgeSeconds,omitempty"`
	}

	LifecycleConfiguration struct {
		Rule []LifecycleRule
	}
	LifecycleRule struct {
		ID         string
		Prefix     string
		Status     string
		Expiration Expiration
	}
	Expiration struct {
		Days int        `xml:"Days,omitempty"`
		Date *time.Time `xml:"Date,omitempty"`
	}

	// BucketLoggingStatus  is the container for logging status information
	BucketLoggingStatus struct {
		LoggingEnabled LoggingEnabled
	}
	// LoggingEnabled is the container for logging information
	LoggingEnabled struct {
		TargetBucket string
		TargetPrefix string
	}

	// WebsiteConfiguration is the container for static website configuration
	WebsiteConfiguration struct {
		IndexDocument IndexDocument
		ErrorDocument ErrorDocument
	}
	// IndexDocument is the container for the Suffix element
	IndexDocument struct {
		Suffix string
	}
	// ErrorDocument is the container for Key element No
	ErrorDocument struct {
		Key string
	}
)

func BucketLocation(value string) Option {
	return xmlBody(&CreateBucketConfiguration{LocationConstraint: value})
}

func newDelete(objects []string, quiet bool) *Delete {
	del := &Delete{Quiet: quiet}
	for _, key := range objects {
		del.Object = append(del.Object, Object{Key: key})
	}
	return del
}

func xmlBody(obj interface{}) Option {
	return func(req *http.Request) error {
		var w bytes.Buffer
		if !pythonSDKCompatibleMode {
			w.WriteString(xml.Header)
		}
		if err := xml.NewEncoder(&w).Encode(obj); err != nil {
			return err
		}
		req.ContentLength = int64(w.Len())
		req.Body = ioutil.NopCloser(bytes.NewReader(w.Bytes()))
		return nil
	}
}

func httpBody(body io.Reader) Option {
	return func(req *http.Request) error {
		rc, ok := body.(io.ReadCloser)
		if !ok && body != nil {
			rc = ioutil.NopCloser(body)
		}
		req.Body = rc
		fileName := ""
		switch v := body.(type) {
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
		case *bytes.Reader:
			req.ContentLength = int64(v.Len())
		case *strings.Reader:
			req.ContentLength = int64(v.Len())
		case *os.File:
			fInfo, err := os.Stat(v.Name())
			if err != nil {
				return err
			}
			fileName = v.Name()
			req.ContentLength = fInfo.Size()
		}
		req.Header.Set("Content-Type", typeByExtension(fileName))
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
