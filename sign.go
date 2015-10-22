package oss

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type signature struct {
	req *http.Request
}

func (s *signature) setContentMD5() {
	if _, ok := s.req.Header["Content-Md5"]; ok {
		return
	}
	buf := []byte{'\n'}
	if s.req.Body != nil {
		buf, _ = ioutil.ReadAll(s.req.Body)
	}
	s.req.Body = ioutil.NopCloser(bytes.NewReader(buf))
	sum := md5.Sum(buf)
	b64 := base64.StdEncoding.EncodeToString(sum[:])
	s.req.Header.Set("Content-Md5", b64)
}

func (s *signature) canonicalizedOSSHeaders() []byte {
	var kvs kvSlice
	for key, vs := range s.req.Header {
		key = strings.ToLower(key)
		if strings.HasPrefix(key, "x-oss-") {
			for _, val := range vs {
				kvs = append(kvs, kv{key, val})
			}
		}
	}
	sort.Sort(kvs)
	var buf bytes.Buffer
	for i, kv := range kvs {
		buf.WriteString(kv.key)
		buf.WriteByte(':')
		buf.WriteString(kv.val)
		if i != len(kvs)-1 {
			buf.WriteByte('\n')
		}
	}
	return buf.Bytes()
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
