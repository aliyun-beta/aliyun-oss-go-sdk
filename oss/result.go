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

type ListAllMyBucketsResult struct {
	Owner   Owner
	Buckets []Bucket `xml:"Buckets>Bucket"`
}

type Owner struct {
	ID          string
	DisplayName string
}

type Bucket struct {
	Location     string
	Name         string
	CreationDate time.Time
}

type ListBucketResult struct {
	Name           string
	Prefix         string
	Marker         string
	MaxKeys        int
	Delimiter      string
	IsTruncated    bool
	Contents       []Content
	CommonPrefixes []string `xml:"CommonPrefixes>Prefix"`
}
type Content struct {
	Key          string
	LastModified time.Time
	ETag         string
	Type         string
	Size         int
	StorageClass string
	Owner        Owner
}

type AccessControlPolicy struct {
	Owner             Owner
	AccessControlList AccessControlList
}

type AccessControlList struct {
	Grant string
}

type LocationConstraint struct {
	Value string `xml:",chardata"`
}

type AppendPosition int

type Header map[string][]string

type DeleteResult struct {
	Deleted []Deleted
}
type Deleted struct {
	Key string
}

type CopyObjectResult struct {
	LastModified string
	ETag         string
}

type InitiateMultipartUploadResult struct {
	Bucket   string
	Key      string
	UploadID string `xml:"UploadId"`
}

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

type writerResult struct {
	io.Writer
}

func (r *writerResult) parse(resp *http.Response) error {
	_, err := io.Copy(r.Writer, resp.Body)
	return err
}

type writeCloserResult struct {
	io.WriteCloser
}

func (r *writeCloserResult) parse(resp *http.Response) error {
	defer r.WriteCloser.Close()
	_, err := io.Copy(r.WriteCloser, resp.Body)
	return err
}
