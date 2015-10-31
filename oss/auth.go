package oss

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type authorization struct {
	req    *http.Request
	secret []byte
}

func setContentMD5(req *http.Request) error {
	if _, ok := req.Header["Content-Md5"]; ok {
		return errors.New("Content-Md5 is already set")
	}
	if req.Body == nil {
		return errors.New("body is nil")
	}
	buf, _ := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewReader(buf))
	sum := md5.Sum(buf)
	b64 := base64.StdEncoding.EncodeToString(sum[:])
	req.Header.Set("Content-Md5", b64)
	return nil
}

func (a *authorization) canonicalizedOSSHeaders() []byte {
	var kvs kvSlice
	for key, vs := range a.req.Header {
		key = strings.ToLower(key)
		if strings.HasPrefix(key, "x-oss-") {
			for _, val := range vs {
				kvs = append(kvs, kv{key, val})
			}
		}
	}
	sort.Sort(kvs)
	var buf bytes.Buffer
	for _, kv := range kvs {
		buf.WriteString(kv.key)
		buf.WriteByte(':')
		buf.WriteString(kv.val)
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

func (a *authorization) data() []byte {
	var w bytes.Buffer
	w.WriteString(a.req.Method)
	w.WriteByte('\n')
	w.WriteString(a.req.Header.Get("Content-Md5"))
	w.WriteByte('\n')
	w.WriteString(a.req.Header.Get("Content-Type"))
	w.WriteByte('\n')
	w.WriteString(a.req.Header.Get("Date"))
	w.WriteByte('\n')
	w.Write(a.canonicalizedOSSHeaders())
	w.WriteString(a.canonicalizedResource())
	return w.Bytes()
}

func (a *authorization) canonicalizedResource() string {
	uri := *a.req.URL
	uri.Host = ""
	uri.Scheme = ""
	return uri.String()
}

func (a *authorization) value() string {
	h := hmac.New(sha1.New, a.secret)
	h.Write(a.data())
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

type (
	kv struct {
		key string
		val string
	}
	kvSlice []kv
)

func (s kvSlice) Len() int           { return len(s) }
func (s kvSlice) Less(i, j int) bool { return s[i].key < s[j].key }
func (s kvSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
