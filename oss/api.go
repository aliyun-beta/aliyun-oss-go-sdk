package oss

import (
	"net/http"
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
	}
}

func (a *API) CreateBucket(name string, acl ACL) error {
	verb := "PUT"
	req, err := http.NewRequest(verb, "http://"+a.endPoint+"/"+name+"/", nil)
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Date", a.now().UTC().Format(gmtTime))
	req.Header.Set("X-Oss-Acl", string(acl))
	req.Header.Set("Authorization", "OSS "+a.accessKeyID+":")
	req.Header.Set("User-Agent", userAgent)
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	return err
}

type Bucket struct {
}
