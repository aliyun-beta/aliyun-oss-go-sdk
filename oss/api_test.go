package oss

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var (
	testTime = "Wed, 21 Oct 2015 15:56:35 GMT"
)

func testNow() time.Time {
	tm, _ := time.Parse(gmtTime, testTime)
	return tm
}

func TestGetService(t *testing.T) {
	rec, err := NewMockServer(`HTTP/1.1 200 OK
Date: Thu, 15 May 2014 11:18:32 GMT
Content-Type: application/xml
Content-Length: 556
Connection: close
Server: AliyunOSS
x-oss-request-id: 5374A2880232A65C23002D74

<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult>
  <Owner>
    <ID>ut_test_put_bucket</ID>
    <DisplayName>ut_test_put_bucket</DisplayName>
  </Owner>
  <Buckets>
    <Bucket>
      <Location>oss-cn-hangzhou-a</Location>
      <Name>xz02tphky6fjfiuc0</Name>
      <CreationDate>2014-05-15T11:18:32.001Z</CreationDate>
    </Bucket>
    <Bucket>
      <Location>oss-cn-hangzhou-a</Location>
      <Name>xz02tphky6fjfiuc1</Name>
      <CreationDate>2014-05-15T11:18:32.002Z</CreationDate>
    </Bucket>
  </Buckets>
</ListAllMyBucketsResult>`)
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	id, secret := "ayahghai0juiSie", "quitie*ph3Lah{F"
	api := New(rec.URL(), id, secret)
	api.now = testNow
	result, err := api.GetService()
	if err != nil {
		t.Fatal(err)
	}
	expected := fmt.Sprintf(`GET / HTTP/1.1
Host: %s
User-Agent: aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:uATT/A0hkaO68KDsD79n17pNA5c=
Date: %s`, rec.URL(), testTime)
	if rec.Err != nil {
		t.Fatal(err)
	}
	if rec.Request != expected {
		t.Fatalf(expectBut, expected, rec.Request)
	}
	expectedResult := &ListAllMyBucketsResult{
		Owner: Owner{
			ID:          "ut_test_put_bucket",
			DisplayName: "ut_test_put_bucket",
		},
		Buckets: []Bucket{
			{
				Location:     "oss-cn-hangzhou-a",
				Name:         "xz02tphky6fjfiuc0",
				CreationDate: parseTime(time.RFC3339Nano, "2014-05-15T11:18:32.001Z"),
			},
			{
				Location:     "oss-cn-hangzhou-a",
				Name:         "xz02tphky6fjfiuc1",
				CreationDate: parseTime(time.RFC3339Nano, "2014-05-15T11:18:32.002Z"),
			},
		},
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf(expectBut, expectedResult, result)
	}
}

func TestPutBucket(t *testing.T) {
	rec, err := NewMockServer("HTTP/1.1 200 OK\n")
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	api := New(rec.URL(), "your-id", "your-secret")
	api.now = testNow
	if err := api.PutBucket("bucket_name", PrivateACL); err != nil {
		t.Fatal(err)
	}
	expected := fmt.Sprintf(`PUT /bucket_name/ HTTP/1.1
Host: %s
User-Agent: aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
Content-Length: 0
Accept-Encoding: identity
Authorization: OSS your-id:h4f1mblhxCOYBJ3jrO+ofuOLO8o=
Date: %s
X-Oss-Acl: private`, rec.URL(), testTime)
	if rec.Err != nil {
		t.Fatal(err)
	}
	if rec.Request != expected {
		t.Fatalf(expectBut, expected, rec.Request)
	}
}

func TestPutObjectFromFile(t *testing.T) {
	rec, err := NewMockServer("HTTP/1.1 200 OK\n")
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	id, secret := "ayahghai0juiSie", "quitie*ph3Lah{F"
	api := New(rec.URL(), id, secret)
	api.now = testNow
	go api.PutObjectFromFile("bucket_name", "object_name", "testdata/test.txt")
	rec.Wait()
	expected := fmt.Sprintf(`PUT /bucket_name/object_name HTTP/1.1
Host: %s
User-Agent: aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
Content-Length: 17
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:CzZjudedDKEx9AkBtLCTnjfuPgE=
Content-Type: text/plain
Date: %s

sfweruewpinbeewa`, rec.URL(), testTime)
	if rec.Err != nil {
		t.Fatal(err)
	}
	if rec.Request != expected {
		t.Fatalf(expectBut, expected, rec.Request)
	}
}

func TestGetBucket(t *testing.T) {
	rec, err := NewMockServer("HTTP/1.1 200 OK\n")
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	id, secret := "ayahghai0juiSie", "quitie*ph3Lah{F"
	api := New(rec.URL(), id, secret)
	api.now = testNow
	go api.GetBucket("bucket_name")
	rec.Wait()
	expected := fmt.Sprintf(`GET /bucket_name/ HTTP/1.1
Host: %s
User-Agent: aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:7hm0/FE4TpkY8OSFunFmTg1TR0Y=
Date: %s`, rec.URL(), testTime)
	if rec.Err != nil {
		t.Fatal(err)
	}
	if rec.Request != expected {
		t.Fatalf(expectBut, expected, rec.Request)
	}
}

func TestGetObjectToFile(t *testing.T) {
	rec, err := NewMockServer("HTTP/1.1 200 OK\n")
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	id, secret := "ayahghai0juiSie", "quitie*ph3Lah{F"
	api := New(rec.URL(), id, secret)
	api.now = testNow
	go api.GetObjectToFile("bucket_name", "object_name", "file.txt")
	rec.Wait()
	expected := fmt.Sprintf(`GET /bucket_name/object_name HTTP/1.1
Host: %s
User-Agent: aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:gkiiIu1xWb5BtqGgF4Pb52mHJWs=
Date: %s`, rec.URL(), testTime)
	if rec.Err != nil {
		t.Fatal(err)
	}
	if rec.Request != expected {
		t.Fatalf(expectBut, expected, rec.Request)
	}
}
