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
	if err != nil {
		return err
	}
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Date", a.now().UTC().Format(gmtTime))
	req.Header.Set("X-Oss-Acl", string(acl))
	req.Header.Set("User-Agent", userAgent)
	auth := authorization{req: req, secret: []byte(a.accessKeySecret)}
	req.Header.Set("Authorization", "OSS "+a.accessKeyID+":"+auth.value())
	_, err = http.DefaultClient.Do(req)
	return err
}

type Bucket struct {
}
