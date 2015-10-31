package oss

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
)

type API struct {
	endPoint        string
	accessKeyID     string
	accessKeySecret string
	now             func() time.Time
	client          *http.Client
}

func New(endPoint, accessKeyID, accessKeySecret string) *API {
	return &API{
		endPoint:        endPoint,
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		now:             time.Now,
		client:          http.DefaultClient,
	}
}

func (a *API) GetService() (res *ListAllMyBucketsResult, _ error) {
	return res, a.do("GET", "", &res)
}

func (a *API) PutBucket(name string, acl ACLType) error {
	return a.do("PUT", name+"/", nil, ACL(acl))
}

func (a *API) GetBucket(name string, options ...Option) (res *ListBucketResult, _ error) {
	return res, a.do("GET", name+"/", &res, options...)
}

func (a *API) GetBucketACL(name string) (res *AccessControlPolicy, _ error) {
	return res, a.do("GET", name+"/?acl", &res)
}

func (a *API) GetBucketLocation(name string) (res *LocationConstraint, _ error) {
	return res, a.do("GET", name+"/?location", &res)
}

func (a *API) DeleteBucket(name string) error {
	return a.do("DELETE", name+"/", nil)
}

func (a *API) GetObjectToFile(bucket, object, fileName string, options ...Option) error {
	// w, err := os.Create(file)
	// if err != nil {
	// 	return err
	// }
	// defer w.Close()
	return a.do("GET", bucket+"/"+object, file(fileName))
}

func (a *API) PutObjectFromString(bucket, object, str string, options ...Option) error {
	return a.do("PUT", bucket+"/"+object, nil,
		append([]Option{ContentType("application/octet-stream"), httpBody(strings.NewReader(str))}, options...)...)
}

func (a *API) PutObjectFromFile(bucket, object, file string, options ...Option) error {
	rd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer rd.Close()
	return a.do("PUT", bucket+"/"+object, nil, append([]Option{httpBody(rd)}, options...)...)
}

func (a *API) AppendObjectFromFile(bucket, object, file string, position AppendPosition, options ...Option) (res AppendPosition, _ error) {
	rd, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer rd.Close()
	return res, a.do("POST", fmt.Sprintf("%s/%s?append&position=%d", bucket, object, position), &res, append([]Option{httpBody(rd)}, options...)...)
}

func (a *API) HeadObject(bucket, object string) (res Header, _ error) {
	return res, a.do("HEAD", bucket+"/"+object, &res)
}

func (a *API) DeleteObject(bucket, object string) error {
	return a.do("DELETE", bucket+"/"+object, nil)
}

func (a *API) do(method, resource string, result interface{}, options ...Option) error {
	req, err := a.newRequest(method, resource, options)
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

func (a *API) newRequest(method, resource string, options []Option) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s/%s", a.endPoint, resource), nil)
	if err != nil {
		return nil, err
	}
	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Date", a.now().UTC().Format(gmtTime))
	req.Header.Set("User-Agent", userAgent)
	auth := authorization{req: req, secret: []byte(a.accessKeySecret)}
	req.Header.Set("Authorization", "OSS "+a.accessKeyID+":"+auth.value())
	return req, nil
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
	if respParser, ok := result.(responseParser); ok {
		return respParser.parse(resp)
	}
	panic(fmt.Sprintf("result %#v should implement responseParser interface", result))
}
