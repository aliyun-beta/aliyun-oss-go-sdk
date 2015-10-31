package oss

import (
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type responseParser interface {
	parse(resp *http.Response) error
}

/*
	if w, ok := result.(io.Writer); ok {
		_, err = io.Copy(w, resp.Body)
		return err
	} else if result != nil {
		return xml.NewDecoder(resp.Body).Decode(result)
	}
*/

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

type file string

func (f file) parse(resp *http.Response) error {
	w, err := os.Create(string(f))
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, resp.Body)
	return err
}
