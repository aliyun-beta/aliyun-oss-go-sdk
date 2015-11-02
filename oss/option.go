package oss

import (
	"net"
	"net/http"
	"strconv"
	"strings"
)

type (
	// Option for an http.Request
	Option func(*http.Request) error
)

var (
	ossHosts = []string{
		"aliyun-inc.com", "aliyuncs.com", "alibaba.net", "s3.amazonaws.com",
	}
)

func bucketHost(bucket string) Option {
	return func(req *http.Request) error {
		if isOSSDomain(req.Host) {
			req.Host = bucket + "." + req.Host
		}
		return nil
	}
}
func isOSSDomain(hostPort string) bool {
	host, _, err := net.SplitHostPort(hostPort)
	if err != nil {
		host = hostPort
	}
	for _, ossHost := range ossHosts {
		if strings.HasSuffix(host, ossHost) {
			return true
		}
	}
	return false
	// ip := net.ParseIP(host)
	// return ip == nil
}

// ACL is an option to set X-Oss-Acl header
func ACL(acl ACLType) Option {
	return setHeader("X-Oss-Acl", string(acl))
}

// ContentType is an option to set Content-Type header
func ContentType(value string) Option {
	return setHeader("Content-Type", value)
}

// CacheControl is an option to set Cache-Control header
func CacheControl(value string) Option {
	return setHeader("Cache-Control", value)
}

// ContentDisposition is an option to set Content-Disposition header
func ContentDisposition(value string) Option {
	return setHeader("Content-Disposition", value)
}

// ContentEncoding is an option to set Content-Encoding header
func ContentEncoding(value string) Option {
	return setHeader("Content-Encoding", value)
}

// Expires is an option to set Expires header
func Expires(value string) Option {
	return setHeader("Expires", value)
}

// Meta is an option to set Meta header
func Meta(key, value string) Option {
	return setHeader("X-Oss-Meta-"+key, value)
}

// Range is an option to set Range header
func Range(value string) Option {
	return setHeader("Range", value)
}

// IfModifiedSince is an option to set If-Modified-Since header
func IfModifiedSince(value string) Option {
	return setHeader("If-Modified-Since", value)
}

// IfUnmodifiedSince is an option to set If-Unmodified-Since header
func IfUnmodifiedSince(value string) Option {
	return setHeader("If-Unmodified-Since", value)
}

// IfMatch is an option to set If-Match header
func IfMatch(value string) Option {
	return setHeader("If-Match", value)
}

// IfNoneMatch is an option to set IfNoneMatch header
func IfNoneMatch(value string) Option {
	return setHeader("If-None-Match", value)
}

// CopySource is an option to set X-Oss-Copy-Source header
func CopySource(sourceBucket, sourceObject string) Option {
	return setHeader("X-Oss-Copy-Source", "/"+sourceBucket+"/"+sourceObject)
}

// CopySourceIfMatch is an option to set X-Oss-Copy-Source-If-Match header
func CopySourceIfMatch(value string) Option {
	return setHeader("X-Oss-Copy-Source-If-Match", value)
}

// CopySourceIfNoneMatch is an option to set X-Oss-Copy-Source-If-None-Match header
func CopySourceIfNoneMatch(value string) Option {
	return setHeader("X-Oss-Copy-Source-If-None-Match", value)
}

// CopySourceIfModifiedSince is an option to set
// X-Oss-CopySource-If-Modified-Since header
func CopySourceIfModifiedSince(value string) Option {
	return setHeader("X-Oss-Copy-Source-If-Modified-Since", value)
}

// CopySourceIfUnmodifiedSince is an option to set
// X-Oss-Copy-Source-If-Unmodified-Since header
func CopySourceIfUnmodifiedSince(value string) Option {
	return setHeader("X-Oss-Copy-Source-If-Unmodified-Since", value)
}

// MetadataDirective is an option to set X-Oss-Metadata-Directive header
func MetadataDirective(directive MetadataDirectiveType) Option {
	return setHeader("X-Oss-Metadata-Directive", string(directive))
}

// ServerSideEncryption is an option to set X-Oss-Server-Side-Encryption header
func ServerSideEncryption(value string) Option {
	return setHeader("X-Oss-Server-Side-Encryption", value)
}

// ObjectACL is an option to set X-Oss-Object-Acl header
func ObjectACL(acl ACLType) Option {
	return setHeader("X-Oss-Object-Acl", string(acl))
}

// ContentLength is an option to set Content-Length header
func ContentLength(length int64) Option {
	return func(req *http.Request) error {
		req.ContentLength = length
		return nil
	}
}
func setHeader(key, value string) Option {
	return func(req *http.Request) error {
		req.Header.Set(key, value)
		return nil
	}
}

// Delimiter is an option to set delimiler parameter
func Delimiter(value string) Option {
	return addParam("delimiter", value)
}

// Marker is an option to set marker parameter
func Marker(value string) Option {
	return addParam("marker", value)
}

// MaxKeys is an option to set maxkeys parameter
func MaxKeys(value int) Option {
	return addParam("maxkeys", strconv.Itoa(value))
}

// Prefix is an option to set prefix parameter
func Prefix(value string) Option {
	return addParam("prefix", value)
}

// EncodingType is an option to set encoding-type parameter
func EncodingType(value string) Option {
	return addParam("encoding-type", value)
}

// ResponseContentType is an option to set response-content-type parameter
func ResponseContentType(value string) Option {
	return addParam("response-content-type", value)
}

// ResponseContentLanguage is an option to set response-content-language parameter
func ResponseContentLanguage(value string) Option {
	return addParam("response-content-language", value)
}

// ResponseCacheControl is an option to set response-cache-control parameter
func ResponseCacheControl(value string) Option {
	return addParam("response-cache-control", value)
}

// ResponseContentDisposition is an option to set response-content-disposition
// parameter
func ResponseContentDisposition(value string) Option {
	return addParam("response-content-disposition", value)
}

// ResponseContentEncoding is an option to set response-content-encoding
// parameter
func ResponseContentEncoding(value string) Option {
	return addParam("response-content-encoding", value)
}

// ResponseExpires is an option to set response-expires parameter
func ResponseExpires(value string) Option {
	return addParam("response-expires", value)
}
func addParam(key, value string) Option {
	return func(req *http.Request) error {
		q := req.URL.Query()
		q.Add(key, value)
		req.URL.RawQuery = q.Encode()
		return nil
	}
}
