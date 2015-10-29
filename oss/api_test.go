package oss

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

const (
	testTimeText = "Wed, 21 Oct 2015 15:56:35 GMT"
	testID       = "ayahghai0juiSie"
	testSecret   = "quitie*ph3Lah{F"
)

var (
	testTime = func() time.Time {
		t, _ := time.Parse(gmtTime, testTimeText)
		return t
	}
)

type testcase struct {
	request          func(*API) (interface{}, error)
	expectedRequest  string
	response         string
	expectedResponse interface{}
	expectedError    error
}

var apiTestcases = []testcase{

	{
		request: func(a *API) (interface{}, error) {
			r, err := a.GetService()
			return r, err
		},
		expectedRequest: `GET / HTTP/1.1
Host: %s
User-Agent: aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:uATT/A0hkaO68KDsD79n17pNA5c=
Date: %s`,
		response: `HTTP/1.1 403 Forbidden

<?xml version="1.0" ?>
<Error xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
    <Code>AccessDenied</Code>
    <Message>Query-string authentication requires the Signature, Expires and OSSAccessKeyId parameters</Message>
    <RequestId>1D842BC5425544BB</RequestId>
    <HostId>oss-cn-hangzhou.aliyuncs.com</HostId>
</Error>`,
		expectedResponse: nil,
		expectedError: &ErrorXML{
			Code:      "AccessDenied",
			Message:   "Query-string authentication requires the Signature, Expires and OSSAccessKeyId parameters",
			RequestID: "1D842BC5425544BB",
			HostID:    "oss-cn-hangzhou.aliyuncs.com",
		},
	},

	{
		request: func(a *API) (interface{}, error) {
			r, err := a.GetService()
			return r, err
		},
		expectedRequest: `GET / HTTP/1.1
Host: %s
User-Agent: aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:uATT/A0hkaO68KDsD79n17pNA5c=
Date: %s`,
		response: `HTTP/1.1 200 OK
Date: Thu, 15 May 2014 11:18:32 GMT
Content-Type: application/xml
Content-Length: 555
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
</ListAllMyBucketsResult>`,
		expectedResponse: &ListAllMyBucketsResult{
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
		},
	},

	{
		request: func(a *API) (interface{}, error) {
			return nil, a.PutBucket("bucket_name", PrivateACL)
		},
		expectedRequest: `PUT /bucket_name/ HTTP/1.1
Host: %s
User-Agent: aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
Content-Length: 0
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:fPC3Cfuif1iGi5LKNRg033EGZcU=
Date: %s
X-Oss-Acl: private`,
		response:         "HTTP/1.1 200 OK\n",
		expectedResponse: nil,
	},

	{
		request: func(a *API) (interface{}, error) {
			r, err := a.GetBucket("bucket_name")
			return r, err
		},
		expectedRequest: `GET /bucket_name/ HTTP/1.1
Host: %s
User-Agent: aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:7hm0/FE4TpkY8OSFunFmTg1TR0Y=
Date: %s`,
		response: `HTTP/1.1 200 OK
x-oss-request-id: 0b05f9b1-539e-a858-0a81-9ca13d8a8011
Date: Fri, 24 Feb 2012 08:43:27 GMT
Content-Type: application/xml
Content-Length: 763
Connection: close
Server: AliyunOSS

<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
<Name>oss-example</Name>
<Prefix>fun/</Prefix>
<Marker></Marker>
<MaxKeys>100</MaxKeys>
<Delimiter>/</Delimiter>
    <IsTruncated>false</IsTruncated>
    <Contents>
        <Key>fun/test.jpg</Key>
        <LastModified>2012-02-24T08:42:32.000Z</LastModified>
        <ETag>&quot;5B3C1A2E053D763E1B002CC607C5A0FE&quot;</ETag>
        <Type>Normal</Type>
        <Size>344606</Size>
        <StorageClass>Standard</StorageClass>
        <Owner>
            <ID>00220120222</ID>
            <DisplayName>user_example</DisplayName>
        </Owner>
    </Contents>
   <CommonPrefixes>
        <Prefix>fun/movie/</Prefix>
   </CommonPrefixes>
</ListBucketResult>`,
		expectedResponse: &ListBucketResult{
			Name:        "oss-example",
			Prefix:      "fun/",
			Marker:      "",
			MaxKeys:     100,
			Delimiter:   "/",
			IsTruncated: false,
			Contents: []Content{
				{
					Key:          "fun/test.jpg",
					LastModified: parseTime(time.RFC3339Nano, "2012-02-24T08:42:32.000Z"),
					ETag:         `"5B3C1A2E053D763E1B002CC607C5A0FE"`,
					Type:         "Normal",
					Size:         344606,
					StorageClass: "Standard",
					Owner: Owner{
						ID:          "00220120222",
						DisplayName: "user_example",
					},
				},
			},
			CommonPrefixes: []string{
				"fun/movie/",
			},
		},
	},

	{
		request: func(a *API) (interface{}, error) {
			return nil, a.PutObjectFromFile("bucket_name", "object_name", "testdata/test.txt")
		},
		expectedRequest: `PUT /bucket_name/object_name HTTP/1.1
Host: %s
User-Agent: aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
Content-Length: 17
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:PrCp4qvn6aHefdgTPfNyW83zPAY=
Content-Type: text/plain; charset=utf-8
Date: %s

sfweruewpinbeewa`,
		response:         "HTTP/1.1 200 OK\n",
		expectedResponse: nil,
	},
}

func TestGetObjectToFile(t *testing.T) {
	rec, err := NewMockServer("HTTP/1.1 200 OK\n")
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	id, secret := "ayahghai0juiSie", "quitie*ph3Lah{F"
	api := New(rec.URL(), id, secret)
	api.now = testTime
	go api.GetObjectToFile("bucket_name", "object_name", "file.txt")
	rec.Wait()
	expected := fmt.Sprintf(`GET /bucket_name/object_name HTTP/1.1
Host: %s
User-Agent: aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:gkiiIu1xWb5BtqGgF4Pb52mHJWs=
Date: %s`, rec.URL(), testTimeText)
	if rec.Err != nil {
		t.Fatal(err)
	}
	if rec.Request != expected {
		t.Fatalf(expectBut, expected, rec.Request)
	}
}

func TestAPIs(t *testing.T) {
	for i := range apiTestcases {
		testAPI(t, i, &apiTestcases[i])
	}
}
func testAPI(t *testing.T, i int, testcase *testcase) {
	rec, err := NewMockServer(testcase.response)
	if err != nil {
		t.Fatal(i, err)
	}
	defer rec.Close()
	api := New(rec.URL(), testID, testSecret)
	api.now = testTime
	response, err := testcase.request(api)
	if v := reflect.ValueOf(response); !v.IsValid() || v.IsNil() {
		response = nil
	}
	if !reflect.DeepEqual(err, testcase.expectedError) {
		t.Fatalf(testcaseExpectBut, i, testcase.expectedError, err)
	}
	expectedRequest := fmt.Sprintf(testcase.expectedRequest, rec.URL(), testTimeText)
	if rec.Err != nil {
		t.Fatal("testcase", i, err)
	}
	if rec.Request != expectedRequest {
		t.Fatalf(testcaseExpectBut, i, expectedRequest, rec.Request)
	}
	if !reflect.DeepEqual(response, testcase.expectedResponse) {
		t.Fatalf(testcaseExpectBut, i, testcase.expectedResponse, response)
	}
}
