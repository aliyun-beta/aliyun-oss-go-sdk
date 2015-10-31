package oss

// Version of the Go SDK
const Version = "0.1.1"

type ACLType string

const (
	PrivateACL      = ACLType("private")
	PublicReadACL   = ACLType("public-read")
	PublicReadWrite = ACLType("public-read-write")
)

const (
	gmtTime = "Mon, 02 Jan 2006 15:04:05 GMT"
)
