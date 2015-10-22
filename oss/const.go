package oss

// Version of the Go SDK
const Version = "0.1.1"

type ACL string

const (
	PrivateACL      = ACL("private")
	PublicReadACL   = ACL("public-read")
	PublicReadWrite = ACL("public-read-write")
)

const (
	gmtTime = "Mon, 02 Jan 2006 15:04:05 GMT"
)
