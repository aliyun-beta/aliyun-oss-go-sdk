package oss

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
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

func (a *API) GetService() error {
	req, err := http.NewRequest("GET", "http://"+a.endPoint+"/", nil)
	if err != nil {
		return err
	}
	resp, err := a.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (a *API) PutBucket(name string, acl ACL) error {
	req, err := http.NewRequest("PUT", "http://"+a.endPoint+"/"+name+"/", nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Oss-Acl", string(acl))
	resp, err := a.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (a *API) GetBucket(name string) error {
	req, err := http.NewRequest("GET", "http://"+a.endPoint+"/"+name+"/", nil)
	if err != nil {
		return err
	}
	resp, err := a.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (a *API) GetObjectToFile(bucket, object, file string) error {
	req, err := http.NewRequest("GET", "http://"+a.endPoint+"/"+bucket+"/"+object, nil)
	if err != nil {
		return err
	}
	resp, err := a.do(req)
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
	fInfo, err := os.Lstat(file)
	if err != nil {
		return err
	}
	rd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer rd.Close()
	req, err := http.NewRequest("PUT", "http://"+a.endPoint+"/"+bucket+"/"+object, rd)
	if err != nil {
		return err
	}
	req.ContentLength = fInfo.Size()
	req.Header.Set("Content-Type", strings.TrimSuffix(mime.TypeByExtension(path.Ext(file)), "; charset=utf-8"))
	resp, err := a.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (a *API) do(req *http.Request) (*http.Response, error) {
	a.setCommonHeaders(req)
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
func (a *API) setCommonHeaders(req *http.Request) {
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Date", a.now().UTC().Format(gmtTime))
	req.Header.Set("User-Agent", userAgent)
	auth := authorization{req: req, secret: []byte(a.accessKeySecret)}
	req.Header.Set("Authorization", "OSS "+a.accessKeyID+":"+auth.value())
}
