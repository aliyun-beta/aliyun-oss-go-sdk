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
	createBucketConfiguration struct {
		XMLName            xml.Name `xml:"CreateBucketConfiguration"`
		LocationConstraint string
	}

	objectsToDelete struct {
		XMLName xml.Name `xml:"Delete"`
		Quiet   bool
		Object  []objectToDelete
	}
	objectToDelete struct {
		Key string
	}

	// CompleteMultipartUpload is the input for CompleteUpload API
	CompleteMultipartUpload struct {
		Part []Part
	}

	// CORSConfiguration represents the CORS rules of a bucket
	CORSConfiguration struct {
		CORSRule []CORSRule
	}
	//CORSRule represents a CORS rule
	CORSRule struct {
		AllowedOrigin []string
		AllowedMethod []string
		AllowedHeader []string
		ExposeHeader  []string
		MaxAgeSeconds int `xml:"MaxAgeSeconds,omitempty"`
	}

	// LifecycleConfiguration represents the deletion rule for a bucket
	LifecycleConfiguration struct {
		Rule []LifecycleRule
	}
	// LifecycleRule represents a deletion rule
	LifecycleRule struct {
		ID         string
		Prefix     string
		Status     string
		Expiration Expiration
	}
	// Expiration represents the expiration time either by days or exact time
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

	// RefererConfiguration is the container for referer configuration
	RefererConfiguration struct {
		AllowEmptyReferer bool
		Referer           []string `xml:"RefererList>Referer"`
	}
)

// BucketLocation is the option for setting bucket location when calling PutBucket
func BucketLocation(value string) Option {
	return XMLBody(&createBucketConfiguration{LocationConstraint: value})
}

func newDelete(objects []string, quiet bool) *objectsToDelete {
	del := &objectsToDelete{Quiet: quiet}
	for _, key := range objects {
		del.Object = append(del.Object, objectToDelete{Key: key})
	}
	return del
}

// XMLBody sets http.Request.Body with XML marshaled from an object
func XMLBody(obj interface{}) Option {
	return func(req *http.Request) error {
		var w bytes.Buffer
		if err := xml.NewEncoder(&w).Encode(obj); err != nil {
			return err
		}
		req.ContentLength = int64(w.Len())
		req.Body = ioutil.NopCloser(bytes.NewReader(w.Bytes()))
		return nil
	}
}

// HTTPBody sets http.Request.Body and Content-Length/Type when possible
func HTTPBody(body io.Reader) Option {
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
			req.ContentLength = tryGetFileSize(v)
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
func tryGetFileSize(f *os.File) int64 {
	fInfo, _ := f.Stat()
	return fInfo.Size()
}
