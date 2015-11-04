package oss

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type (
	// API is the entry object for all OSS methods
	API struct {
		endPoint        string
		accessKeyID     string
		accessKeySecret string
		securityToken   string
		client          *http.Client
		scheme          string
		now             func() time.Time
	}
	// APIOption provides optional configurations for an API object
	APIOption func(*API)
)

// New creates an API object
func New(endPoint, accessKeyID, accessKeySecret string, options ...APIOption) *API {
	api := &API{
		endPoint:        endPoint,
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		client:          http.DefaultClient,
		scheme:          "http",
		now:             time.Now,
	}
	for _, option := range options {
		option(api)
	}
	return api
}

// HTTPClient sets the underlying http.Client object used by the OSS client.
func HTTPClient(client *http.Client) APIOption {
	return func(a *API) {
		if client != nil {
			a.client = client
		}
	}
}

// SecurityToken sets the STS token for temporary access
func SecurityToken(token string) APIOption {
	return func(a *API) {
		a.securityToken = token
	}
}

// URLScheme sets the scheme for the API access: http or https, default is http.
func URLScheme(scheme string) APIOption {
	return func(a *API) {
		switch scheme {
		case "http", "https":
			a.scheme = scheme
		}
	}
}

// GetService list all buckets
func (a *API) GetService(options ...Option) (res *ListAllMyBucketsResult, _ error) {
	return res, a.Do("GET", "", "", &res, options...)
}

// PutBucket creates a new bucket
func (a *API) PutBucket(name string, acl ACLType) error {
	return a.Do("PUT", name, "", nil, ACL(acl))
}

// PutBucketACL sets acess right for a bucket
func (a *API) PutBucketACL(name string, acl ACLType) error {
	return a.Do("PUT", name, "?acl", nil, ACL(acl))
}

// PutBucketLogging configures a bucket's logging behavior
func (a *API) PutBucketLogging(name string, status *BucketLoggingStatus) error {
	return a.Do("PUT", name, "?logging", nil, XMLBody(status))
}

// PutBucketWebsite configures a bucket as a static website
func (a *API) PutBucketWebsite(name string, config *WebsiteConfiguration) error {
	return a.Do("PUT", name, "?website", nil, XMLBody(config))
}

// PutBucketReferer configures a bucket's referer whitelist
func (a *API) PutBucketReferer(name string, config *RefererConfiguration) error {
	return a.Do("PUT", name, "?referer", nil, XMLBody(config))
}

// PutBucketLifecycle configures the automatic deletion of a bucket
func (a *API) PutBucketLifecycle(bucket string, lifecycle *LifecycleConfiguration) error {
	return a.Do("PUT", bucket, "?lifecycle", nil, XMLBody(lifecycle))
}

// GetBucket returns all the objects in a bucket
func (a *API) GetBucket(name string, options ...Option) (res *ListBucketResult, _ error) {
	return res, a.Do("GET", name, "", &res, options...)
}

// GetBucketACL returns the access rule for a bucket
func (a *API) GetBucketACL(name string) (res *AccessControlPolicy, _ error) {
	return res, a.Do("GET", name, "?acl", &res)
}

// GetBucketLocation returns the location of a bucket
func (a *API) GetBucketLocation(name string) (res *LocationConstraint, _ error) {
	return res, a.Do("GET", name, "?location", &res)
}

// GetBucketLogging returns a bucket's logging configuration
func (a *API) GetBucketLogging(name string) (res *BucketLoggingStatus, _ error) {
	return res, a.Do("GET", name, "?logging", &res)
}

// GetBucketWebsite returns a bucket's configuration as a static website
func (a *API) GetBucketWebsite(name string) (res *WebsiteConfiguration, _ error) {
	return res, a.Do("GET", name, "?website", &res)
}

// GetBucketReferer returns a bucket's referer whitelist
func (a *API) GetBucketReferer(name string) (res *RefererConfiguration, _ error) {
	return res, a.Do("GET", name, "?referer", &res)
}

// GetBucketLifecycle returns a bucket's deletion configuration
func (a *API) GetBucketLifecycle(bucket string) (res *LifecycleConfiguration, _ error) {
	return res, a.Do("GET", bucket, "?lifecycle", &res)
}

// DeleteBucket deletes a bucket
func (a *API) DeleteBucket(name string) error {
	return a.Do("DELETE", name, "", nil)
}

// DeleteBucketLogging turns off the logging functionality
func (a *API) DeleteBucketLogging(name string) error {
	return a.Do("DELETE", name, "?logging", nil)
}

// DeleteBucketWebsite turns off the website functionality
func (a *API) DeleteBucketWebsite(name string) error {
	return a.Do("DELETE", name, "?website", nil)
}

// DeleteBucketLifecycle deletes the lifecycle configuration of a bucket
func (a *API) DeleteBucketLifecycle(bucket string) error {
	return a.Do("DELETE", bucket, "?lifecycle", nil)
}

// PutObject uploads a file from an io.Reader
func (a *API) PutObject(bucket, object string, rd io.Reader, options ...Option) error {
	return a.Do("PUT", bucket, object, nil, append([]Option{HTTPBody(rd)}, options...)...)
}

// CopyObject copies an existing object on OSS to another object
func (a *API) CopyObject(sourceBucket, sourceObject, targetBucket, targetObject string, options ...Option) (res *CopyObjectResult, _ error) {
	return res, a.Do("PUT", targetBucket, targetObject, &res, append(options, CopySource(sourceBucket, sourceObject))...)
}

// GetObject returns an object and write it to an io.Writer
func (a *API) GetObject(bucket, object string, w io.Writer, options ...Option) (res Header, _ error) {
	return res, a.Do("GET", bucket, object, &bodyAndHeader{Writer: w, Header: &res})
}

// AppendObject uploads a file by append to it from an io.Reader
func (a *API) AppendObject(bucket, object string, rd io.Reader, position AppendPosition, options ...Option) (res AppendPosition, _ error) {
	return res, a.Do("POST", bucket, fmt.Sprintf("%s?append&position=%d", object, position), &res, append([]Option{HTTPBody(rd)}, options...)...)
}

// DeleteObject deletes an object
func (a *API) DeleteObject(bucket, object string) error {
	return a.Do("DELETE", bucket, object, nil)
}

// DeleteObjects deletes multiple objects
func (a *API) DeleteObjects(bucket string, quiet bool, objects ...string) (res *DeleteResult, _ error) {
	return res, a.Do("POST", bucket, "?delete", &res, XMLBody(newDelete(objects, quiet)), ContentMD5)
}

// HeadObject returns only the metadata of an object in HTTP headers
func (a *API) HeadObject(bucket, object string) (res Header, _ error) {
	return res, a.Do("HEAD", bucket, object, &res)
}

// PutObjectACL sets acess right for an object
func (a *API) PutObjectACL(bucket, object string, acl ACLType) error {
	return a.Do("PUT", bucket, object+"?acl", nil, ACL(acl))
}

// GetObjectACL returns the access rule for an object
func (a *API) GetObjectACL(bucket, object string) (res *AccessControlPolicy, _ error) {
	return res, a.Do("GET", bucket, object+"?acl", &res)
}

// InitUpload starts an multipart upload process
func (a *API) InitUpload(bucket, object string, options ...Option) (res *InitiateMultipartUploadResult, _ error) {
	return res, a.Do("POST", bucket, object+"?uploads", &res, append(options, ContentType("application/octet-stream"))...)
}

// UploadPart updates a trunk of data from an io.Reader
func (a *API) UploadPart(bucket, object string, uploadID string, partNumber int, rd io.Reader, size int64) (res *UploadPartResult, _ error) {
	return res, a.Do("PUT", bucket, fmt.Sprintf("%s?partNumber=%d&uploadId=%s", object, partNumber, uploadID), &res, HTTPBody(&io.LimitedReader{R: rd, N: size}), ContentLength(size))
}

// UploadPartCopy updates a trunk of data from an existing object
func (a *API) UploadPartCopy(bucket, object string, uploadID string, partNumber int, sourceBucket, sourceObject string, options ...Option) (res *CopyPartResult, _ error) {
	return res, a.Do("PUT", bucket, fmt.Sprintf("%s?partNumber=%d&uploadId=%s", object, partNumber, uploadID), &res, CopySource(sourceBucket, sourceObject))
}

// CompleteUpload notifies that the multipart upload is complete
func (a *API) CompleteUpload(bucket, object string, uploadID string, list *CompleteMultipartUpload) (res *CompleteMultipartUploadResult, _ error) {
	return res, a.Do("POST", bucket, fmt.Sprintf("%s?uploadId=%s", object, uploadID), &res, XMLBody(list), ContentMD5, ContentType("application/octet-stream"))
}

// AbortUpload aborts a multipart upload
func (a *API) AbortUpload(bucket, object string, uploadID string) error {
	return a.Do("DELETE", bucket, fmt.Sprintf("%s?uploadId=%s", object, uploadID), nil)
}

// ListUploads lists all ongoing multipart uploads
func (a *API) ListUploads(bucket, object string, options ...Option) (res *ListMultipartUploadsResult, _ error) {
	return res, a.Do("GET", bucket, "?uploads", &res, options...)
}

// ListParts lists successful uploaded parts of a multipart upload
func (a *API) ListParts(bucket, object, uploadID string, options ...Option) (res *ListPartsResult, _ error) {
	return res, a.Do("GET", bucket, fmt.Sprintf("%s?uploadId=%s", object, uploadID), &res, options...)
}

// PutBucketCORS sets CORS rules to a bucket
func (a *API) PutBucketCORS(bucket string, cors *CORSConfiguration) error {
	return a.Do("PUT", bucket, "?cors", nil, XMLBody(cors), ContentMD5)
}

// GetBucketCORS gets CORS rules of a bucket
func (a *API) GetBucketCORS(bucket string) (res *CORSConfiguration, _ error) {
	return res, a.Do("GET", bucket, "?cors", &res)
}

// DeleteBucketCORS deletes the CORS rules of a bucket
func (a *API) DeleteBucketCORS(bucket string) error {
	return a.Do("DELETE", bucket, "?cors", nil)
}

// OptionObject queries OSS whether a CORS request is permitted or not
func (a *API) OptionObject(bucket, object string, options ...Option) (res Header, _ error) {
	return res, a.Do("OPTIONS", bucket, object, &res, options...)
}

// Do sends a general OSS request and returns the response.
func (a *API) Do(method, bucket, object string, result interface{}, options ...Option) error {
	req, err := a.newRequest(method, bucket, object, options)
	if err != nil {
		return err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return a.handleResponse(resp, result)
}

func (a *API) newRequest(method, bucket, object string, options []Option) (*http.Request, error) {
	uri, err := ossURL(a.scheme, a.endPoint, bucket, object)
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		Method:     method,
		URL:        uri,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       uri.Host,
	}
	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Date", a.now().UTC().Format(gmtTime))
	req.Header.Set("User-Agent", userAgent)
	if a.securityToken != "" {
		req.Header.Set("X-Oss-Security-Token", a.securityToken)
	}
	if !strings.HasPrefix(req.Header.Get("Content-Type"), "multipart/form-data") {
		auth := authorization{
			req:    req,
			bucket: bucket,
			object: object,
			secret: []byte(a.accessKeySecret),
		}
		req.Header.Set("Authorization", "OSS "+a.accessKeyID+":"+auth.value())
	}
	return req, nil
}

var (
	rxBucketName = regexp.MustCompile(`\A[a-z0-9][a-z0-9\-]{2,62}\z`)
	rxObjectName = regexp.MustCompile(`\A[^/\\](|[^\r\n]*)\z`)
)

func ossURL(scheme, host, bucket, object string) (*url.URL, error) {
	scheme += "://"
	if bucket == "" && object == "" {
		return url.Parse(scheme + host)
	}
	if !rxBucketName.MatchString(bucket) {
		return nil, ErrInvalidBucketName
	}
	if len(object) > 1023 || object != "" && !rxObjectName.MatchString(object) {
		return nil, ErrInvalidObjectName
	}
	if !isOSSDomain(host) {
		return url.Parse(scheme + fmt.Sprintf("%s/%s/%s", host, bucket, object))
	}
	return url.Parse(scheme + fmt.Sprintf("%s.%s/%s", bucket, host, object))
}

func (a *API) handleResponse(resp *http.Response, result interface{}) error {
	if resp.StatusCode/100 > 2 {
		return parseError(resp)
	}
	if result == nil {
		return nil
	}
	if v := reflect.ValueOf(result); v.Kind() == reflect.Ptr {
		if v.Elem().Kind() == reflect.Ptr {
			v = v.Elem()
			v.Set(reflect.New(v.Type().Elem()))
			result = v.Interface()
		}
	}
	if respParser, ok := result.(ResponseParser); ok {
		return respParser.Parse(resp)
	}
	panic(fmt.Sprintf("programming error: type of %#v should implement responseParser interface", result))
}
