package oss

import (
	"fmt"
	"testing"
	"time"
)

func TestPutBucket(t *testing.T) {
	rec, err := NewRequestRecorder()
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	api := New(rec.URL(), "your-id", "your-secret")
	now := "Wed, 21 Oct 2015 15:56:35 GMT"
	api.now = func() time.Time {
		tm, err := time.Parse(gmtTime, now)
		if err != nil {
			t.Fatal(err)
		}
		return tm
	}
	go api.CreateBucket("bucket_name", PrivateACL)
	rec.Wait()
	expected := fmt.Sprintf(`PUT /bucket_name/ HTTP/1.1
Host: %s
User-Agent: aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
Content-Length: 0
Accept-Encoding: identity
Authorization: OSS your-id:qKkbehZghgryl+VGYs0rE83gV2g=
Date: %s
X-Oss-Acl: private`, rec.URL(), now)
	if rec.Err != nil {
		t.Fatal(err)
	}
	if rec.Request != expected {
		t.Fatalf("\n--- EXPECT ---\n%s\n--- BUT GOT ---\n%s", expected, rec.Request)
	}
}
