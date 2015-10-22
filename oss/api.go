package oss

import (
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

func (a *API) PutBucket(name string, acl ACL) error {
	verb := "PUT"
	req, err := http.NewRequest(verb, "http://"+a.endPoint+"/"+name+"/", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Date", a.now().UTC().Format(gmtTime))
	req.Header.Set("X-Oss-Acl", string(acl))
	req.Header.Set("User-Agent", userAgent)
	auth := authorization{req: req, secret: []byte(a.accessKeySecret)}
	req.Header.Set("Authorization", "OSS "+a.accessKeyID+":"+auth.value())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return err
}

func (a *API) PutObjectFromFile(bucket, object, file string) error {
	verb := "PUT"
	fInfo, err := os.Lstat(file)
	if err != nil {
		return err
	}
	rd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer rd.Close()
	req, err := http.NewRequest(verb, "http://"+a.endPoint+"/"+bucket+"/"+object, rd)
	if err != nil {
		return err
	}
	req.ContentLength = fInfo.Size()
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Date", a.now().UTC().Format(gmtTime))
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", strings.TrimSuffix(mime.TypeByExtension(path.Ext(file)), "; charset=utf-8"))
	auth := authorization{req: req, secret: []byte(a.accessKeySecret)}
	req.Header.Set("Authorization", "OSS "+a.accessKeyID+":"+auth.value())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return err
}

type Bucket struct {
}
