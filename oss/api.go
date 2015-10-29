package oss

import (
	"encoding/xml"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
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

func (a *API) GetService() (*ListAllMyBucketsResult, error) {
	resp, err := a.do("GET", "", nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result := new(ListAllMyBucketsResult)
	return result, xml.NewDecoder(resp.Body).Decode(result)
}

func (a *API) PutBucket(name string, acl ACL) error {
	resp, err := a.do("PUT", name+"/", http.Header{"X-Oss-Acl": []string{string(acl)}}, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (a *API) GetBucket(name string) (*ListBucketResult, error) {
	resp, err := a.do("GET", name+"/", nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result := new(ListBucketResult)
	return result, xml.NewDecoder(resp.Body).Decode(result)
}

func (a *API) GetObjectToFile(bucket, object, file string) error {
	resp, err := a.do("GET", bucket+"/"+object, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	w, err := os.Create(file)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, resp.Body)
	return err
}

func (a *API) PutObjectFromFile(bucket, object, file string) error {
	rd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer rd.Close()
	resp, err := a.do("PUT", bucket+"/"+object, nil, rd)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (a *API) do(method, resource string, header http.Header, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s/%s", a.endPoint, resource), body)
	if err != nil {
		return nil, err
	}
	if header != nil {
		req.Header = header
	}
	if err := a.setCommonHeaders(req); err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 > 2 {
		defer resp.Body.Close()
		if xmlError, err := parseErrorXML(resp.Body); err != nil {
			return nil, err
		} else {
			return nil, xmlError
		}
	}
	return resp, nil
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
