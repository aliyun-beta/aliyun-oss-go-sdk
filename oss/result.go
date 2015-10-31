package oss

import (
	"net/http"
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
