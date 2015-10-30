package oss

import (
	"encoding/xml"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"reflect"
	"time"
)

type API struct {
	endPoint        string
	accessKeyID     string
	accessKeySecret string
	now             func() time.Time
}

func New(endPoint, accessKeyID, accessKeySecret string) *API {
	return &API{
		endPoint:        endPoint,
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		now:             time.Now,
	}
}

func (a *API) GetService() (res *ListAllMyBucketsResult, _ error) {
	return res, a.do("GET", "", nil, nil, &res)
}

func (a *API) PutBucket(name string, acl ACL) error {
	return a.do("PUT", name+"/", http.Header{"X-Oss-Acl": []string{string(acl)}}, nil, nil)
}

func (a *API) GetBucket(name string) (res *ListBucketResult, _ error) {
	return res, a.do("GET", name+"/", nil, nil, &res)
}

func (a *API) GetObjectToFile(bucket, object, file string) error {
	w, err := os.Create(file)
	if err != nil {
		return err
	}
	defer w.Close()
	return a.do("GET", bucket+"/"+object, nil, nil, w)
}

func (a *API) PutObjectFromFile(bucket, object, file string) error {
	rd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer rd.Close()
	return a.do("PUT", bucket+"/"+object, nil, rd, nil)
}

func (a *API) do(method, resource string, header http.Header, body io.Reader, result interface{}) error {
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s/%s", a.endPoint, resource), body)
	if err != nil {
		return err
	}
	if header != nil {
		req.Header = header
	}
	if err := a.setCommonHeaders(req); err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 > 2 {
		if xmlError, err := parseErrorXML(resp.Body); err != nil {
			return err
		} else {
			return xmlError
		}
	}
	if w, ok := result.(io.Writer); ok {
		_, err = io.Copy(w, resp.Body)
		return err
	} else if result != nil {
		alloc(reflect.ValueOf(result))
		return xml.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}
func alloc(v reflect.Value) {
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
}
func (a *API) setCommonHeaders(req *http.Request) error {
	if f, ok := req.Body.(*os.File); ok {
		fInfo, err := os.Stat(f.Name())
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", mime.TypeByExtension(path.Ext(f.Name())))
		req.ContentLength = fInfo.Size()
	}
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Date", a.now().UTC().Format(gmtTime))
	req.Header.Set("User-Agent", userAgent)
	auth := authorization{req: req, secret: []byte(a.accessKeySecret)}
	req.Header.Set("Authorization", "OSS "+a.accessKeyID+":"+auth.value())
	return nil
}
