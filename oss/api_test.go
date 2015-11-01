package oss

import (
	"fmt"
	"reflect"
	"strings"
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
	name             string
	request          func(*API) (interface{}, error)
	expectedRequest  string
	response         string
	expectedResponse interface{}
	expectedError    error
}

var apiTestcases = []testcase{

	{
		name: "authorization fail",
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
		name: "GetService",
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
		name: "PutBucket",
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
		name: "GetBucket",
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
		name: "PutObjectFromFile",
		request: func(a *API) (interface{}, error) {
			return nil, a.PutObjectFromFile("bucket_name", "object_name", testFileName)
		},
		expectedRequest: `PUT /bucket_name/object_name HTTP/1.1
Host: %s
User-Agent: %s
Content-Length: 16
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:cUN83rKdXAq2MRbzQZYWJC4hIRg=
Content-Type: application/octet-stream
Date: %s

sfweruewpinbeewa`,
		response:         "HTTP/1.1 200 OK\n",
		expectedResponse: nil,
	},

	{
		name: "PutObjectFromString",
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
		name: "GetBucketACL",
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
		name: "GetBucketLocation",
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
		name: "DeleteBucket",
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
		name: "AppendObjectFromFile",
		request: func(a *API) (interface{}, error) {
			r, err := a.AppendObjectFromFile(testBucketName, testObjectName, testFileName, 0)
			return r, err
		},
		expectedRequest: `POST /bucket_name/object_name?append&position=0 HTTP/1.1
Host: %s
User-Agent: %s
Content-Length: 16
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
x-oss-next-append-position: 16
x-oss-request-id: 559CC9BDC755F95A64485981
`,
		expectedResponse: AppendPosition(16),
	},

	{
		name: "HeadObject",
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
		name: "DeleteObject",
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
		name: "DeleteObjects",
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
		name: "CopyObject",
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

	{
		name: "InitUpload",
		request: func(a *API) (interface{}, error) {
			r, err := a.InitUpload(testBucketName, testObjectName)
			return r, err
		},
		expectedRequest: `POST /bucket_name/object_name?uploads HTTP/1.1
Host: %s
User-Agent: %s
Content-Length: 0
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:i0ZMmOYrcPZibjOyuazIcpqP45Q=
Content-Type: application/octet-stream
Date: %s`,
		response: `HTTP/1.1 200 OK
Content-Length: 273
Server: AliyunOSS
Connection: close
x-oss-request-id: 42c25703-7503-fbd8-670a-bda01eaec618
Date: Wed, 22 Feb 2012 08:32:21 GMT
Content-Type: application/xml

<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
    <Bucket>bucket_name</Bucket>
    <Key>multipart.data</Key>
    <UploadId>0004B9894A22E5B1888A1E29F8236E2D</UploadId>
</InitiateMultipartUploadResult>`,
		expectedResponse: &InitiateMultipartUploadResult{
			Bucket:   "bucket_name",
			Key:      "multipart.data",
			UploadID: "0004B9894A22E5B1888A1E29F8236E2D",
		},
	},

	{
		name: "UploadPart",
		request: func(a *API) (interface{}, error) {
			r, err := a.UploadPart(testBucketName, testObjectName, "0004B9895DBBB6EC98E36", 1, strings.NewReader(`sfweruewpinbeewa`), 16)
			return r, err
		},
		expectedRequest: `PUT /bucket_name/object_name?partNumber=1&uploadId=0004B9895DBBB6EC98E36 HTTP/1.1
Host: %s
User-Agent: %s
Content-Length: 16
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:im5VQx0x3baFKjafCmEiDxClVU0=
Content-Type: application/octet-stream
Date: %s

sfweruewpinbeewa`,
		response: `HTTP/1.1 200 OK
Server: AliyunOSS
Connection: close
ETag: 7265F4D211B56873A381D321F586E4A9
x-oss-request-id: 3e6aba62-1eae-d246-6118-8ff42cd0c21a
Date: Wed, 22 Feb 2012 08:32:21 GMT
`,
		expectedResponse: &UploadPartResult{ETag: "7265F4D211B56873A381D321F586E4A9"},
	},

	{
		name: "CompleteUpload",
		request: func(a *API) (interface{}, error) {
			list := &CompleteMultipartUpload{
				Part: []Part{
					{
						PartNumber: 1,
						ETag:       `C1B61751512FFC8B0E86675D114497A6`,
					},
				},
			}
			r, err := a.CompleteUpload(testBucketName, testObjectName, "0004B9895DBBB6EC98E36", list)
			return r, err
		},
		expectedRequest: `POST /bucket_name/object_name?uploadId=0004B9895DBBB6EC98E36 HTTP/1.1
Host: %s
User-Agent: %s
Content-Length: 174
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:vnZp7SnuGTXV5KPYq+GV+ws5Ff4=
Content-Md5: zJBKGLrHC8XqtxxS7kZM+Q==
Date: %s

<?xml version="1.0" encoding="UTF-8"?>
<CompleteMultipartUpload><Part><PartNumber>1</PartNumber><ETag>C1B61751512FFC8B0E86675D114497A6</ETag></Part></CompleteMultipartUpload>`,
		response: `HTTP/1.1 200 OK
Server: AliyunOSS
Content-Length: 356
Content-Type: Application/xml
Connection: close
x-oss-request-id: 594f0751-3b1e-168f-4501-4ac71d217d6e
Date: Fri, 24 Feb 2012 10:19:18 GMT

<?xml version="1.0" encoding="UTF-8"?>
<CompleteMultipartUploadResult xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
    <Location>http://oss-example.oss-cn-hangzhou.aliyuncs.com /multipart.data</Location>
    <Bucket>oss-example</Bucket>
    <Key>multipart.data</Key>
    <ETag>B864DB6A936D376F9F8D3ED3BBE540DD-3</ETag>
</CompleteMultipartUploadResult>`,
		expectedResponse: &CompleteMultipartUploadResult{
			Location: "http://oss-example.oss-cn-hangzhou.aliyuncs.com /multipart.data",
			Bucket:   "oss-example",
			Key:      "multipart.data",
			ETag:     "B864DB6A936D376F9F8D3ED3BBE540DD-3",
		},
	},

	{
		name: "CancelUpload",
		request: func(a *API) (interface{}, error) {
			return nil, a.CancelUpload(testBucketName, testObjectName, "0004B9895DBBB6EC98E36")
		},
		expectedRequest: `DELETE /bucket_name/object_name?uploadId=0004B9895DBBB6EC98E36 HTTP/1.1
Host: %s
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:afD0ln6LveomIhMndS4klSzfCwM=
Date: %s`,
		response: `HTTP/1.1 204
Server: AliyunOSS
Connection: close
x-oss-request-id: 059a22ba-6ba9-daed-5f3a-e48027df344d
Date: Wed, 22 Feb 2012 08:32:21 GMT
`,
		expectedResponse: nil,
	},

	{
		name: "ListUploads",
		request: func(a *API) (interface{}, error) {
			r, err := a.ListUploads(testBucketName, testObjectName)
			return r, err
		},
		expectedRequest: `GET /bucket_name/?uploads HTTP/1.1
Host: %s
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:EJSXM8cV0QJNk9jPFGgaeARRu5Y=
Date: %s`,
		response: `HTTP/1.1 200
Server: AliyunOSS
Connection: close
Content-length: 1839
Content-type: application/xml
x-oss-request-id: 58a41847-3d93-1905-20db-ba6f561ce67a
Date: Thu, 23 Feb 2012 06:14:27 GMT

<?xml version="1.0" encoding="UTF-8"?>
<ListMultipartUploadsResult xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
    <Bucket>oss-example</Bucket>
    <KeyMarker></KeyMarker>
    <UploadIdMarker></UploadIdMarker>
    <NextKeyMarker>oss.avi</NextKeyMarker>
    <NextUploadIdMarker>0004B99B8E707874FC2D692FA5D77D3F</NextUploadIdMarker>
    <Delimiter></Delimiter>
    <Prefix></Prefix>
    <MaxUploads>1000</MaxUploads>
    <IsTruncated>false</IsTruncated>
    <Upload>
        <Key>multipart.data</Key>
        <UploadId>0004B999EF518A1FE585B0C9360DC4C8</UploadId>
        <Initiated>2012-02-23T04:18:23.000Z</Initiated>
    </Upload>
    <Upload>
        <Key>multipart.data</Key>
        <UploadId>0004B999EF5A239BB9138C6227D69F95</UploadId>
        <Initiated>2012-02-23T04:18:23.000Z</Initiated>
    </Upload>
    <Upload>
        <Key>oss.avi</Key>
        <UploadId>0004B99B8E707874FC2D692FA5D77D3F</UploadId>
        <Initiated>2012-02-23T06:14:27.000Z</Initiated>
    </Upload>
</ListMultipartUploadsResult>`,
		expectedResponse: &ListMultipartUploadsResult{
			Bucket:             "oss-example",
			NextKeyMarker:      "oss.avi",
			NextUploadIDMarker: "0004B99B8E707874FC2D692FA5D77D3F",
			MaxUploads:         1000,
			IsTruncated:        false,
			Upload: []Upload{
				{
					Key:       "multipart.data",
					UploadID:  "0004B999EF518A1FE585B0C9360DC4C8",
					Initiated: parseTime(time.RFC3339Nano, "2012-02-23T04:18:23.000Z"),
				},
				{
					Key:       "multipart.data",
					UploadID:  "0004B999EF5A239BB9138C6227D69F95",
					Initiated: parseTime(time.RFC3339Nano, "2012-02-23T04:18:23.000Z"),
				},
				{
					Key:       "oss.avi",
					UploadID:  "0004B99B8E707874FC2D692FA5D77D3F",
					Initiated: parseTime(time.RFC3339Nano, "2012-02-23T06:14:27.000Z"),
				},
			},
		},
	},

	{
		name: "ListParts",
		request: func(a *API) (interface{}, error) {
			r, err := a.ListParts(testBucketName, testObjectName, "0004B9895DBBB6EC98E36")
			return r, err
		},
		expectedRequest: `GET /bucket_name/object_name?uploadId=0004B9895DBBB6EC98E36 HTTP/1.1
Host: %s
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:/zw9uMvuU3TtkYmnFM1SVoXI6P8=
Date: %s`,
		response: `HTTP/1.1 200
Server: AliyunOSS
Connection: close
Content-length: 1221
Content-type: application/xml
x-oss-request-id: 106452c8-10ff-812d-736e-c865294afc1c
Date: Thu, 23 Feb 2012 07:13:28 GMT

<?xml version="1.0" encoding="UTF-8"?>
<ListPartsResult xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
    <Bucket>multipart_upload</Bucket>
    <Key>multipart.data</Key>
    <UploadId>0004B999EF5A239BB9138C6227D69F95</UploadId>
    <NextPartNumberMarker>5</NextPartNumberMarker>
    <MaxParts>1000</MaxParts>
    <IsTruncated>false</IsTruncated>
    <Part>
        <PartNumber>1</PartNumber>
        <LastModified>2012-02-23T07:01:34.000Z</LastModified>
        <ETag>&quot;3349DC700140D7F86A078484278075A9&quot;</ETag>
        <Size>6291456</Size>
    </Part>
    <Part>
        <PartNumber>2</PartNumber>
        <LastModified>2012-02-23T07:01:12.000Z</LastModified>
        <ETag>&quot;3349DC700140D7F86A078484278075A9&quot;</ETag>
        <Size>6291456</Size>
    </Part>
    <Part>
        <PartNumber>5</PartNumber>
        <LastModified>2012-02-23T07:02:03.000Z</LastModified>
        <ETag>&quot;7265F4D211B56873A381D321F586E4A9&quot;</ETag>
        <Size>1024</Size>
    </Part>
</ListPartsResult>`,
		expectedResponse: &ListPartsResult{
			Bucket:               "multipart_upload",
			Key:                  "multipart.data",
			UploadID:             "0004B999EF5A239BB9138C6227D69F95",
			NextPartNumberMarker: 5,
			MaxParts:             1000,
			IsTruncated:          false,
			Part: []Part{
				{
					PartNumber:   1,
					LastModified: parseTimePtr(time.RFC3339Nano, "2012-02-23T07:01:34.000Z"),
					ETag:         `"3349DC700140D7F86A078484278075A9"`,
					Size:         6291456,
				},
				{
					PartNumber:   2,
					LastModified: parseTimePtr(time.RFC3339Nano, "2012-02-23T07:01:12.000Z"),
					ETag:         `"3349DC700140D7F86A078484278075A9"`,
					Size:         6291456,
				},
				{
					PartNumber:   5,
					LastModified: parseTimePtr(time.RFC3339Nano, "2012-02-23T07:02:03.000Z"),
					ETag:         `"7265F4D211B56873A381D321F586E4A9"`,
					Size:         1024,
				},
			},
		},
	},

	{
		name: "PutCORS",
		request: func(a *API) (interface{}, error) {
			return nil, a.PutCORS(testBucketName, &CORSConfiguration{
				CORSRule: []CORSRule{
					{
						AllowedOrigin: []string{"*"},
						AllowedMethod: []string{"PUT", "GET"},
						AllowedHeader: []string{"Authorization"},
					},
					{
						AllowedOrigin: []string{"http://www.a.com", "http://www.b.com"},
						AllowedMethod: []string{"GET"},
						AllowedHeader: []string{"Authorization"},
						ExposeHeader:  []string{"x-oss-test", "x-oss-test1"},
						MaxAgeSeconds: 100,
					},
				},
			})
		},
		expectedRequest: `PUT /bucket_name/?cors HTTP/1.1
Host: %s
User-Agent: %s
Content-Length: 549
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:AZ9v6CeAbKAWaphXg3GHuZ26FuM=
Content-Md5: dCyZ6ocGwvoNqax+nPRAlg==
Date: %s

<?xml version="1.0" encoding="UTF-8"?>
<CORSConfiguration><CORSRule><AllowedOrigin>*</AllowedOrigin><AllowedMethod>PUT</AllowedMethod><AllowedMethod>GET</AllowedMethod><AllowedHeader>Authorization</AllowedHeader></CORSRule><CORSRule><AllowedOrigin>http://www.a.com</AllowedOrigin><AllowedOrigin>http://www.b.com</AllowedOrigin><AllowedMethod>GET</AllowedMethod><AllowedHeader>Authorization</AllowedHeader><ExposeHeader>x-oss-test</ExposeHeader><ExposeHeader>x-oss-test1</ExposeHeader><MaxAgeSeconds>100</MaxAgeSeconds></CORSRule></CORSConfiguration>`,
		response: `HTTP/1.1 200 OK
x-oss-request-id: 50519080C4689A033D00235F
Date: Fri, 04 May 2012 03:21:12 GMT
Content-Length: 0
Connection: close
Server: AliyunOSS
`,
		expectedResponse: nil,
	},

	{
		name: "GetCORS",
		request: func(a *API) (interface{}, error) {
			r, err := a.GetCORS(testBucketName)
			return r, err
		},
		expectedRequest: `GET /bucket_name/?cors HTTP/1.1
Host: %s
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:PBEn27rJaD+O4E/2gN22KbaMbcE=
Date: %s`,
		response: `HTTP/1.1 200
x-oss-request-id: 50519080C4689A033D00235F
Date: Thu, 13 Sep 2012 07:51:28 GMT
Connection: close
Content-Length: 317
Server: AliyunOSS

<?xml version="1.0" encoding="UTF-8"?>
<CORSConfiguration>
    <CORSRule>
      <AllowedOrigin>*</AllowedOrigin>
      <AllowedMethod>GET</AllowedMethod>
      <AllowedHeader>*</AllowedHeader>
      <ExposeHeader>x-oss-test</ExposeHeader>
      <MaxAgeSeconds>100</MaxAgeSeconds>
    </CORSRule>
</CORSConfiguration>`,
		expectedResponse: &CORSConfiguration{
			CORSRule: []CORSRule{
				{
					AllowedOrigin: []string{"*"},
					AllowedMethod: []string{"GET"},
					AllowedHeader: []string{"*"},
					ExposeHeader:  []string{"x-oss-test"},
					MaxAgeSeconds: 100,
				},
			},
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
		testAPI(t, &apiTestcases[i])
	}
}
func testAPI(t *testing.T, testcase *testcase) {
	rec, err := NewMockServer(testcase.response)
	if err != nil {
		t.Fatalf(testcaseErr, testcase.name, err)
	}
	defer rec.Close()
	api := New(rec.URL(), testID, testSecret)
	api.now = testTime
	response, err := testcase.request(api)
	if v := reflect.ValueOf(response); !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
		response = nil
	}
	if !reflect.DeepEqual(err, testcase.expectedError) {
		t.Fatalf(testcaseExpectBut, testcase.name, testcase.expectedError, err)
	}
	expectedRequest := fmt.Sprintf(testcase.expectedRequest, rec.URL(), userAgent, testTimeText)
	if rec.Err != nil {
		t.Fatalf(testcaseErr, testcase.name, err)
	}
	if rec.Request != expectedRequest {
		t.Fatalf(testcaseExpectBut, testcase.name, expectedRequest, rec.Request)
	}
	if !reflect.DeepEqual(response, testcase.expectedResponse) {
		t.Fatalf(testcaseExpectBut, testcase.name, testcase.expectedResponse, response)
	}
}
