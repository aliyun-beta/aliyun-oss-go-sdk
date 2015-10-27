package oss

import (
	"time"
)

type ListAllMyBucketsResult struct {
	Owner   Owner    `xml:"Owner"`
	Buckets []Bucket `xml:"Buckets>Bucket"`
}

type Owner struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

type Bucket struct {
	Location     string    `xml:"Location"`
	Name         string    `xml:"Name"`
	CreationDate time.Time `xml:"CreationDate"`
}
