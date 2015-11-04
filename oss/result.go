package oss

import (
	"encoding/xml"
	"io"
	"net/http"
	"strconv"
	"time"
)

// ResponseParser parses an HTTP response to an object
type ResponseParser interface {
	Parse(resp *http.Response) error
}

// XML result types
type (
	// ListAllMyBucketsResult is returned by GetService API
	ListAllMyBucketsResult struct {
		Owner   Owner
		Buckets []Bucket `xml:"Buckets>Bucket"`
	}
	// Owner of a bucket
	Owner struct {
		ID          string
		DisplayName string
	}
	// Bucket information
	Bucket struct {
		Location     string
		Name         string
		CreationDate time.Time
	}

	// ListBucketResult is returned by GetBucket API
	ListBucketResult struct {
		Name           string
		Prefix         string
		Marker         string
		MaxKeys        int
		Delimiter      string
		IsTruncated    bool
		Contents       []Content
		CommonPrefixes []string `xml:"CommonPrefixes>Prefix"`
	}
	// Content is the container of an object's meta information
	Content struct {
		Key          string
		LastModified time.Time
		ETag         string
		Type         string
		Size         int64
		StorageClass string
		Owner        Owner
	}

	// AccessControlPolicy is returned by GetBucketACL API
	AccessControlPolicy struct {
		Owner             Owner
		AccessControlList AccessControlList
	}
	// AccessControlList is the container of ACL information
	AccessControlList struct {
		Grant string
	}

	// LocationConstraint represents the data center of a bucket
	LocationConstraint struct {
		Value string `xml:",chardata"`
	}

	// AppendPosition is returned by AppendObject API
	AppendPosition int

	// Header object contains the HTTP headers
	Header map[string][]string

	// DeleteResult is returned by DeleteObjects API
	DeleteResult struct {
		Deleted []Deleted
	}
	// Deleted is the container of a deleted object key
	Deleted struct {
		Key string
	}

	// CopyObjectResult is returned by CopyObject API
	CopyObjectResult struct {
		LastModified string
		ETag         string
	}

	// InitiateMultipartUploadResult is returned by InitUpload API
	InitiateMultipartUploadResult struct {
		Bucket   string
		Key      string
		UploadID string `xml:"UploadId"`
	}

	// ListMultipartUploadsResult is returned by ListUploads API
	ListMultipartUploadsResult struct {
		Bucket             string
		EncodingType       string
		KeyMarker          string
		UploadIDMarker     string `xml:"UploadIdMarker"`
		NextKeyMarker      string
		NextUploadIDMarker string `xml:"NextUploadIdMarker"`
		Delimiter          string
		Prefix             string
		MaxUploads         int
		IsTruncated        bool
		Upload             []Upload
	}
	// Upload represents the information of an upload object
	Upload struct {
		Key       string
		UploadID  string `xml:"UploadId"`
		Initiated time.Time
	}

	// CopyPartResult is returned by UploadPartCopy API
	CopyPartResult struct {
		LastModified *time.Time
		ETag         string
	}

	// CompleteMultipartUploadResult is returned by CompleteUpload API
	CompleteMultipartUploadResult struct {
		Location string
		Bucket   string
		Key      string
		ETag     string
	}

	// ListPartsResult is returned by ListParts API
	ListPartsResult struct {
		Bucket               string
		EncodingType         string
		Key                  string
		UploadID             string `xml:"UploadId"`
		PartNumberMarker     int
		NextPartNumberMarker int
		MaxParts             int
		IsTruncated          bool
		Part                 []Part
	}
	// Part represents the information of an upload part
	Part struct {
		PartNumber   int
		LastModified *time.Time `xml:"LastModified,omitempty"`
		ETag         string
		Size         int64 `xml:"Size,omitempty"`
	}
)

// Parse implements ResponseParser
func (r *LocationConstraint) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *ListAllMyBucketsResult) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *ListBucketResult) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *AccessControlPolicy) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *AppendPosition) Parse(resp *http.Response) error {
	i, err := strconv.Atoi(resp.Header.Get("X-Oss-Next-Append-Position"))
	if err != nil {
		return err
	}
	*r = AppendPosition(i)
	return nil
}

// Parse implements ResponseParser
func (r *Header) Parse(resp *http.Response) error {
	*r = Header(resp.Header)
	return nil
}

// Parse implements ResponseParser
func (r *DeleteResult) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *CopyObjectResult) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *InitiateMultipartUploadResult) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *CompleteMultipartUploadResult) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *ListMultipartUploadsResult) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *ListPartsResult) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *CORSConfiguration) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *LifecycleConfiguration) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *BucketLoggingStatus) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *WebsiteConfiguration) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *RefererConfiguration) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// Parse implements ResponseParser
func (r *CopyPartResult) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

// bodyAndHeader writes http.Response.Body to an io.Writer and also saves Header object
type bodyAndHeader struct {
	io.Writer
	*Header
}

// Parse implements ResponseParser
func (r *bodyAndHeader) Parse(resp *http.Response) error {
	*r.Header = Header(resp.Header)
	_, err := io.Copy(r.Writer, resp.Body)
	return err
}

// UploadPartResult is the container of the ETag returned by UploadPart API
type UploadPartResult struct {
	ETag string
}

// Parse implements ResponseParser
func (r *UploadPartResult) Parse(resp *http.Response) error {
	r.ETag = resp.Header.Get("ETag")
	return nil
}
