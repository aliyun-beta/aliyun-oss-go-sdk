package oss

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

const (
	testTimeText   = "Wed, 21 Oct 2015 15:56:35 GMT"
	testID         = "ayahghai0juiSie"
	testSecret     = "quitie*ph3Lah{F"
	testBucketName = "bucket_name"
	testObjectName = "object_name"
	testFileName   = "testdata/test"
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

	// authorization fail
	{
		request: func(a *API) (interface{}, error) {
			r, err := a.GetService()
			return r, err
		},
		expectedRequest: `GET / HTTP/1.1
Host: %s
User-Agent: %s
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
		expectedError: &Error{
			Code:           "AccessDenied",
			Message:        "Query-string authentication requires the Signature, Expires and OSSAccessKeyId parameters",
			RequestID:      "1D842BC5425544BB",
			HostID:         "oss-cn-hangzhou.aliyuncs.com",
			HTTPStatusCode: 403,
			HTTPStatus:     "403 Forbidden",
		},
	},

	{
		request: func(a *API) (interface{}, error) {
			r, err := a.GetService()
			return r, err
		},
		expectedRequest: `GET / HTTP/1.1
Host: %s
User-Agent: %s
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
User-Agent: %s
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
User-Agent: %s
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
User-Agent: %s
Content-Length: 17
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:PrCp4qvn6aHefdgTPfNyW83zPAY=
Content-Type: text/plain; charset=utf-8
Date: %s

sfweruewpinbeewa`,
		response:         "HTTP/1.1 200 OK\n",
		expectedResponse: nil,
	},

	{
		request: func(a *API) (interface{}, error) {
			return nil, a.PutObjectFromString("bucket_name", "object_name", "wefpofjwefew")
		},
		expectedRequest: `PUT /bucket_name/object_name HTTP/1.1
Host: %s
User-Agent: %s
Content-Length: 12
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:cUN83rKdXAq2MRbzQZYWJC4hIRg=
Content-Type: application/octet-stream
Date: %s

wefpofjwefew`,
		response:         "HTTP/1.1 200 OK\n",
		expectedResponse: nil,
	},

	{
		request: func(a *API) (interface{}, error) {
			r, err := a.GetBucketACL("bucket_name")
			return r, err
		},
		expectedRequest: `GET /bucket_name/?acl HTTP/1.1
Host: %s
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:LKyy+96PVVFIFgIYgwz3DD3CEzs=
Date: %s`,
		response: `HTTP/1.1 200 OK
x-oss-request-id: 6f720c98-40fe-6de0-047b-e7fb08c4059b
Date: Fri, 24 Feb 2012 04:11:23 GMT
Content-Length: 253
Content-Tupe: application/xml
Connection: close
Server: AliyunOSS

<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>00220120222</ID>
        <DisplayName>user_example</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>public-read</Grant>
    </AccessControlList>
</AccessControlPolicy>`,
		expectedResponse: &AccessControlPolicy{
			Owner: Owner{
				ID:          "00220120222",
				DisplayName: "user_example",
			},
			AccessControlList: AccessControlList{
				Grant: "public-read",
			},
		},
	},

	{
		request: func(a *API) (interface{}, error) {
			r, err := a.GetBucketLocation(testBucketName)
			return r, err
		},
		expectedRequest: `GET /bucket_name/?location HTTP/1.1
Host: %s
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:ycP0hM0Uk40gkqXhljFeHVTWkko=
Date: %s`,
		response: `HTTP/1.1 200
x-oss-request-id: 513836E0F687780D1A690708
Date: Fri, 15 Mar 2013 05:31:04 GMT
Connection: close
Content-Length: 143
Server: AliyunOSS

<?xml version="1.0" encoding="UTF-8"?>
<LocationConstraint xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">oss-cn-hangzhou</LocationConstraint>`,
		expectedResponse: &LocationConstraint{
			Value: "oss-cn-hangzhou",
		},
	},

	{
		request: func(a *API) (interface{}, error) {
			return nil, a.DeleteBucket(testBucketName)
		},
		expectedRequest: `DELETE /bucket_name/ HTTP/1.1
Host: %s
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:L4mtpy/LSUGq1J/oNTOd/lnMEw8=
Date: %s`,
		response:         "HTTP/1.1 200\n",
		expectedResponse: nil,
	},

	{
		request: func(a *API) (interface{}, error) {
			r, err := a.AppendObjectFromFile(testBucketName, testObjectName, testFileName, 0)
			return r, err
		},
		expectedRequest: `POST /bucket_name/object_name?append&position=0 HTTP/1.1
Host: %s
User-Agent: %s
Content-Length: 17
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:pTwcHEynVLcFuA99DVwYrxr2nlk=
Content-Type: application/octet-stream
Date: %s

sfweruewpinbeewa`,
		response: `HTTP/1.1 200 OK
Date: Wed, 08 Jul 2015 06:57:01 GMT
ETag: "0F7230CAA4BE94CCBDC99C5500000000"
Connection: close
Content-Length: 0
Server: AliyunOSS
x-oss-hash-crc64ecma: 14741617095266562575
x-oss-next-append-position: 17
x-oss-request-id: 559CC9BDC755F95A64485981
`,
		expectedResponse: AppendPosition(17),
	},

	{
		request: func(a *API) (interface{}, error) {
			r, err := a.HeadObject(testBucketName, testObjectName)
			return r, err
		},
		expectedRequest: `HEAD /bucket_name/object_name HTTP/1.1
Host: %s
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:nK5UyJ5cD6+AShftVF8YQI2+Oo4=
Date: %s`,
		response: `HTTP/1.1 200 OK
x-oss-request-id: 06d4be30-2216-9264-757a-8f8b19b254bb
ETag: "fba9dede5f27731c9771645a39863328"
Content-Length: 344606
`,
		expectedResponse: Header{
			"X-Oss-Request-Id": []string{"06d4be30-2216-9264-757a-8f8b19b254bb"},
			"Etag":             []string{`"fba9dede5f27731c9771645a39863328"`},
			"Content-Length":   []string{"344606"},
		},
	},

	{
		request: func(a *API) (interface{}, error) {
			return nil, a.DeleteObject(testBucketName, testObjectName)
		},
		expectedRequest: `DELETE /bucket_name/object_name HTTP/1.1
Host: %s
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:SAfVfVFR6w1tFpfqE0xBGZaryb8=
Date: %s`,
		response:         "HTTP/1.1 200 OK\n",
		expectedResponse: nil,
	},

	{
		request: func(a *API) (interface{}, error) {
			r, err := a.DeleteObjects(testBucketName, false, "obj1", "obj2", "obj3")
			return r, err
		},
		expectedRequest: `POST /bucket_name/?delete HTTP/1.1
Host: %s
User-Agent: %s
Content-Length: 172
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:Jk4IiThihXCj8bmIwFPbc7kHbco=
Content-Md5: Tbx4zkqDSc6oNRnTo4dndg==
Date: %s

<?xml version="1.0" encoding="UTF-8"?>
<Delete><Quiet>false</Quiet><Object><Key>obj1</Key></Object><Object><Key>obj2</Key></Object><Object><Key>obj3</Key></Object></Delete>`,
		response: `HTTP/1.1 200 OK
x-oss-request-id: 78320852-7eee-b697-75e1-b6db0f4849e7
Date: Wed, 29 Feb 2012 12:26:16 GMT
Content-Length: 274
Content-Type: application/xml
Connection: close
Server: AliyunOSS

<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
    <Deleted>
       <Key>obj1</Key>
    </Deleted>
    <Deleted>
       <Key>obj2</Key>
    </Deleted>
    <Deleted>
       <Key>obj3</Key>
    </Deleted>
</DeleteResult>`,
		expectedResponse: &DeleteResult{
			Deleted: []Deleted{
				{
					Key: "obj1",
				},
				{
					Key: "obj2",
				},
				{
					Key: "obj3",
				},
			},
		},
	},

	{
		request: func(a *API) (interface{}, error) {
			r, err := a.CopyObject("source_bucket", "source_object", "target_bucket", "target_object")
			return r, err
		},
		expectedRequest: `PUT /target_bucket/target_object HTTP/1.1
Host: %s
User-Agent: %s
Content-Length: 0
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:mRUja421nH8GsCJJ4vWuQszdW8g=
Date: %s
X-Oss-Copy-Source: /source_bucket/source_object`,
		response: `HTTP/1.1 200 OK
x-oss-request-id: 3dfb2597-72a0-b3f7-320f-8b6627a96e68
Content-Type: application/xml
Content-Length: 241
Connection: close
Date: Fri, 24 Feb 2012 07:18:48 GMT
Server: AliyunOSS

<?xml version="1.0" encoding="UTF-8"?>
<CopyObjectResult xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
    <LastModified>Fri, 24 Feb 2012 07:18:48 GMT</LastModified>
    <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE"</ETag>
</CopyObjectResult>`,
		expectedResponse: &CopyObjectResult{
			LastModified: "Fri, 24 Feb 2012 07:18:48 GMT",
			ETag:         `"5B3C1A2E053D763E1B002CC607C5A0FE"`,
		},
	},
}

func TestGetObjectToFile(t *testing.T) {
	rec, err := NewMockServer(`HTTP/1.1 200 OK
x-oss-request-id: 3a89276f-2e2d-7965-3ff9-51c875b99c41
x-oss-object-type: Normal
Date: Fri, 24 Feb 2012 06:38:30 GMT
Last-Modified: Fri, 24 Feb 2012 06:07:48 GMT
ETag: "5B3C1A2E053D763E1B002CC607C5A0FE "
Content-Type: text/plain
Content-Length: 7
Server: AliyunOSS

abcdef`)
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()
	api := New(rec.URL(), testID, testSecret)
	api.now = testTime
	err = api.GetObjectToFile(testBucketName, testObjectName, "file.txt")
	if err != nil {
		t.Fatal(err)
	}
	expected := fmt.Sprintf(`GET /bucket_name/object_name HTTP/1.1
Host: %s
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:gkiiIu1xWb5BtqGgF4Pb52mHJWs=
Date: %s`, rec.URL(), userAgent, testTimeText)
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
	if v := reflect.ValueOf(response); !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
		response = nil
	}
	if !reflect.DeepEqual(err, testcase.expectedError) {
		t.Fatalf(testcaseExpectBut, i, testcase.expectedError, err)
	}
	expectedRequest := fmt.Sprintf(testcase.expectedRequest, rec.URL(), userAgent, testTimeText)
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
