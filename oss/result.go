package oss

import (
	"encoding/xml"
	"io"
	"net/http"
	"strconv"
	"time"
)

type responseParser interface {
	parse(resp *http.Response) error
}

// XML result types
type (
	ListAllMyBucketsResult struct {
		Owner   Owner
		Buckets []Bucket `xml:"Buckets>Bucket"`
	}
	Owner struct {
		ID          string
		DisplayName string
	}
	Bucket struct {
		Location     string
		Name         string
		CreationDate time.Time
	}

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
	Content struct {
		Key          string
		LastModified time.Time
		ETag         string
		Type         string
		Size         int
		StorageClass string
		Owner        Owner
	}

	AccessControlPolicy struct {
		Owner             Owner
		AccessControlList AccessControlList
	}
	AccessControlList struct {
		Grant string
	}

	LocationConstraint struct {
		Value string `xml:",chardata"`
	}

	AppendPosition int

	Header map[string][]string

	DeleteResult struct {
		Deleted []Deleted
	}
	Deleted struct {
		Key string
	}

	CopyObjectResult struct {
		LastModified string
		ETag         string
	}

	InitiateMultipartUploadResult struct {
		Bucket   string
		Key      string
		UploadID string `xml:"UploadId"`
	}

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
	Upload struct {
		Key       string
		UploadID  string `xml:"UploadId"`
		Initiated time.Time
	}

	CompleteMultipartUploadResult struct {
		Location string
		Bucket   string
		Key      string
		ETag     string
	}

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
	Part struct {
		PartNumber   int
		LastModified *time.Time `xml:"LastModified,omitempty"`
		ETag         string
		Size         int `xml:"Size,omitempty"`
	}
)

func (r *LocationConstraint) parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}
func (r *ListAllMyBucketsResult) parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}
func (r *ListBucketResult) parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}
func (r *AccessControlPolicy) parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}
func (r *AppendPosition) parse(resp *http.Response) error {
	i, err := strconv.Atoi(resp.Header.Get("X-Oss-Next-Append-Position"))
	if err != nil {
		return err
	}
	*r = AppendPosition(i)
	return nil
}
func (r *Header) parse(resp *http.Response) error {
	*r = Header(resp.Header)
	return nil
}
func (r *DeleteResult) parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}
func (r *CopyObjectResult) parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}
func (r *InitiateMultipartUploadResult) parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}
func (r *CompleteMultipartUploadResult) parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}
func (r *ListMultipartUploadsResult) parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}
func (r *ListPartsResult) parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}
func (r *CORSConfiguration) parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}
func (r *LifecycleConfiguration) parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}
func (r *BucketLoggingStatus) parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}

type writerResult struct {
	io.Writer
}

func (r *writerResult) parse(resp *http.Response) error {
	_, err := io.Copy(r.Writer, resp.Body)
	return err
}

type UploadPartResult struct {
	ETag string
}

func (r *UploadPartResult) parse(resp *http.Response) error {
	r.ETag = resp.Header.Get("ETag")
	return nil
}
