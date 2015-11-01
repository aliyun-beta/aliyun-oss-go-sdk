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

func ACL(acl ACLType) Option {
	return setHeader("X-Oss-Acl", string(acl))
}
func ContentType(value string) Option {
	return setHeader("Content-Type", value)
}
func CacheControl(value string) Option {
	return setHeader("Cache-Control", value)
}
func ContentDisposition(value string) Option {
	return setHeader("Content-Disposition", value)
}
func ContentEncoding(value string) Option {
	return setHeader("Content-Encoding", value)
}
func Expires(value string) Option {
	return setHeader("Expires", value)
}
func Meta(key, value string) Option {
	return setHeader("X-Oss-Meta-"+key, value)
}
func Range(value string) Option {
	return setHeader("Range", value)
}
func IfModifiedSince(value string) Option {
	return setHeader("If-Modified-Since", value)
}
func IfUnmodifiedSince(value string) Option {
	return setHeader("If-Unmodified-Since", value)
}
func IfMatch(value string) Option {
	return setHeader("If-Match", value)
}
func IfNoneMatch(value string) Option {
	return setHeader("If-None-Match", value)
}
func CopySource(sourceBucket, sourceObject string) Option {
	return setHeader("X-Oss-Copy-Source", "/"+sourceBucket+"/"+sourceObject)
}
func CopySourceIfMatch(value string) Option {
	return setHeader("X-Oss-Copy-Source-If-Match", value)
}
func CopySourceIfNoneMatch(value string) Option {
	return setHeader("X-Oss-Copy-Source-If-None-Match", value)
}
func CopySourceIfModifiedSince(value string) Option {
	return setHeader("X-Oss-Copy-Source-If-Modified-Since", value)
}
func CopySourceIfUnmodifiedSince(value string) Option {
	return setHeader("X-Oss-Copy-Source-If-Unmodified-Since", value)
}
func MetadataDirective(directive MetadataDirectiveType) Option {
	return setHeader("X-Oss-Metadata-Directive", string(directive))
}
func ServerSideEncryption(value string) Option {
	return setHeader("X-Oss-Server-Side-Encryption", value)
}
func ObjectACL(acl ACLType) Option {
	return setHeader("X-Oss-Object-Acl", string(acl))
}
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

func Delimiter(value string) Option {
	return addParam("delimiter", value)
}
func Marker(value string) Option {
	return addParam("marker", value)
}
func MaxKeys(value int) Option {
	return addParam("maxkeys", strconv.Itoa(value))
}
func Prefix(value string) Option {
	return addParam("prefix", value)
}
func EncodingType(value string) Option {
	return addParam("encoding-type", value)
}
func ResponseContentType(value string) Option {
	return addParam("Response-Content-Type", value)
}
func ResponseContentLanguage(value string) Option {
	return addParam("Response-Content-Language", value)
}
func ResponseCacheControl(value string) Option {
	return addParam("Response-Cache-Control", value)
}
func ResponseContentDisposition(value string) Option {
	return addParam("Response-Content-Disposition", value)
}
func ResponseContentEncoding(value string) Option {
	return addParam("Response-Content-Encoding", value)
}
func ResponseExpires(value string) Option {
	return addParam("Response-Expires", value)
}
func addParam(key, value string) Option {
	return func(req *http.Request) error {
		q := req.URL.Query()
		q.Add(key, value)
		req.URL.RawQuery = q.Encode()
		return nil
	}
}
