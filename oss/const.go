package oss

// Version of the Go SDK
const Version = "0.1.1"

// ACLType contains possible values of Access Control List
type ACLType string

const (
	// PrivateACL represents private access right
	PrivateACL = ACLType("private")
	// PublicReadACL represents public-read acess right
	PublicReadACL = ACLType("public-read")
	// PublicReadWriteACL represents public-read-write acess right
	PublicReadWriteACL = ACLType("public-read-write")
)

// MetadataDirectiveType contains possbile values of metadata directive
type MetadataDirectiveType string

const (
	// CopyMeta represents COPY directive
	CopyMeta = MetadataDirectiveType("COPY")
	// ReplaceMeta represents REPLACE directive
	ReplaceMeta = MetadataDirectiveType("REPLACE")
)

const (
	gmtTime = "Mon, 02 Jan 2006 15:04:05 GMT"
)
