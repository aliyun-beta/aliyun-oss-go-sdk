package oss

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

const (
	testTimeText   = "Wed, 21 Oct 2015 15:56:35 GMT"
	testID         = "ayahghai0juiSie"
	testSecret     = "quitie*ph3Lah{F"
	testBucketName = "bucket-name"
	testObjectName = "object/name"
	testFileName   = "testdata/test"
	testEndpoint   = "oss-cn-hangzhou.aliyuncs.com"
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
Host: oss-cn-hangzhou.aliyuncs.com
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
Host: oss-cn-hangzhou.aliyuncs.com
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
			return nil, a.PutBucket("bucket-name", PrivateACL)
		},
		expectedRequest: `PUT / HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 0
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:g4VuHXgdXYmY+EnBsd6GPpc1fk0=
Date: %s
X-Oss-Acl: private`,
		response:         "HTTP/1.1 200 OK\n",
		expectedResponse: nil,
	},

	{
		name: "PutBucketACL",
		request: func(a *API) (interface{}, error) {
			return nil, a.PutBucketACL("bucket-name", PublicReadACL)
		},
		expectedRequest: `PUT /?acl HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 0
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:+lDnwIxmtWknsG2MUXsMivVpZoQ=
Date: %s
X-Oss-Acl: public-read`,
		response:         "HTTP/1.1 200 OK\n",
		expectedResponse: nil,
	},

	{
		name: "PutBucketLogging",
		request: func(a *API) (interface{}, error) {
			return nil, a.PutBucketLogging("bucket-name", &BucketLoggingStatus{
				LoggingEnabled: LoggingEnabled{
					TargetBucket: "doc-log",
					TargetPrefix: "MyLog-",
				},
			})
		},
		expectedRequest: `PUT /?logging HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 147
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:5rVsETIss359jA9q2k8+RTwk09g=
Date: %s

<BucketLoggingStatus><LoggingEnabled><TargetBucket>doc-log</TargetBucket><TargetPrefix>MyLog-</TargetPrefix></LoggingEnabled></BucketLoggingStatus>`,
		response: `HTTP/1.1 200 OK
x-oss-request-id: 19a86d66-3492-0465-12af-7bec0938d0f9
Date: Fri, 04 May 2012 03:21:12 GMT
Content-Length: 0
Connection: close
Server: AliyunOSS
`,
		expectedResponse: nil,
	},

	{
		name: "PutBucketWebsite",
		request: func(a *API) (interface{}, error) {
			return nil, a.PutBucketWebsite("bucket-name", &WebsiteConfiguration{
				IndexDocument: IndexDocument{
					Suffix: "index.html",
				},
				ErrorDocument: ErrorDocument{
					Key: "error.html",
				},
			})
		},
		expectedRequest: `PUT /?website HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 155
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:cYJRyRil8SL+NTHnjgc37kHDPFI=
Date: %s

<WebsiteConfiguration><IndexDocument><Suffix>index.html</Suffix></IndexDocument><ErrorDocument><Key>error.html</Key></ErrorDocument></WebsiteConfiguration>`,
		response: `HTTP/1.1 200 OK
x-oss-request-id: 19a86d66-3492-0465-12af-7bec0938d0f9
Date: Fri, 04 May 2012 03:21:12 GMT
Content-Length: 0
Connection: close
Server: AliyunOSS
`,
		expectedResponse: nil,
	},

	{
		name: "PutBucketReferer",
		request: func(a *API) (interface{}, error) {
			return nil, a.PutBucketReferer("bucket-name", &RefererConfiguration{
				AllowEmptyReferer: true,
				Referer: []string{
					"http://www.aliyun.com",
					"https://www.aliyun.com",
				},
			})
		},
		expectedRequest: `PUT /?referer HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 196
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:vPGa/wuveVoHH857+l04tOpZn3A=
Date: %s

<RefererConfiguration><AllowEmptyReferer>true</AllowEmptyReferer><RefererList><Referer>http://www.aliyun.com</Referer><Referer>https://www.aliyun.com</Referer></RefererList></RefererConfiguration>`,
		response: `HTTP/1.1 200 OK
x-oss-request-id: 19a86d66-3492-0465-12af-7bec0938d0f9
Date: Fri, 04 May 2012 03:21:12 GMT
Content-Length: 0
Connection: close
Server: AliyunOSS
`,
		expectedResponse: nil,
	},

	{
		name: "PutBucketLifecycle",
		request: func(a *API) (interface{}, error) {
			lifecycle := &LifecycleConfiguration{
				Rule: []LifecycleRule{
					{
						ID:     "delete obsoleted files",
						Prefix: "obsoleted/",
						Status: "Enabled",
						Expiration: Expiration{
							Days: 3,
						},
					},
					{
						ID:     "delete temporary files",
						Prefix: "temporary/",
						Status: "Enabled",
						Expiration: Expiration{
							Date: parseTimePtr(time.RFC3339Nano, "2022-10-12T00:00:00.001Z"),
						},
					},
				},
			}
			return nil, a.PutBucketLifecycle(testBucketName, lifecycle)
		},
		expectedRequest: `PUT /?lifecycle HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 340
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:3jdkrgCFCOELpNpH0K4rxkjK6AU=
Date: %s

<LifecycleConfiguration><Rule><ID>delete obsoleted files</ID><Prefix>obsoleted/</Prefix><Status>Enabled</Status><Expiration><Days>3</Days></Expiration></Rule><Rule><ID>delete temporary files</ID><Prefix>temporary/</Prefix><Status>Enabled</Status><Expiration><Date>2022-10-12T00:00:00.001Z</Date></Expiration></Rule></LifecycleConfiguration>`,
		response: `HTTP/1.1 200 OK
x-oss-request-id: 534B371674E88A4D8906008B
Date: Mon, 14 Apr 2014 01:17:10 GMT
Content-Length: 0
Connection: close
Server: AliyunOSS
`,
		expectedResponse: nil,
	},

	{
		name: "GetBucket",
		request: func(a *API) (interface{}, error) {
			r, err := a.GetBucket(testBucketName)
			return r, err
		},
		expectedRequest: `GET / HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:VVbyxjdp2eJ8g5t7o7XxlFy0kNo=
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
		name: "GetBucketACL",
		request: func(a *API) (interface{}, error) {
			r, err := a.GetBucketACL("bucket-name")
			return r, err
		},
		expectedRequest: `GET /?acl HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:boNBCRWMvPxob0Wkv+5qB9RNGeQ=
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
		expectedRequest: `GET /?location HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:LrbZbg8dvqId7QbH9j8+JjuNAwQ=
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
		expectedRequest: `DELETE / HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:FXJsV//QmBZqMczqQkR2Lm+QUKY=
Date: %s`,
		response:         "HTTP/1.1 200\n",
		expectedResponse: nil,
	},

	{
		name: "GetBucketLogging",
		request: func(a *API) (interface{}, error) {
			r, err := a.GetBucketLogging(testBucketName)
			return r, err
		},
		expectedRequest: `GET /?logging HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:1iVz2H1ae00yFT6f5zf9KNIV2n4=
Date: %s`,
		response: `HTTP/1.1 200
x-oss-request-id: 7faf664d-0cad-852e-4b38-2ac2232e7e7f
Date: Fri, 04 May 2012 05:31:04 GMT
Connection: close
Content-Length: 259
Server: AliyunOSS

<?xml version="1.0" encoding="UTF-8"?>
<BucketLoggingStatus xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
<LoggingEnabled>
<TargetBucket>mybucketlogs</TargetBucket>
<TargetPrefix>mybucket-access_log/</TargetPrefix>
</LoggingEnabled>
</BucketLoggingStatus>`,
		expectedResponse: &BucketLoggingStatus{
			LoggingEnabled: LoggingEnabled{
				TargetBucket: "mybucketlogs",
				TargetPrefix: "mybucket-access_log/",
			},
		},
	},

	{
		name: "GetBucketWebsite",
		request: func(a *API) (interface{}, error) {
			r, err := a.GetBucketWebsite(testBucketName)
			return r, err
		},
		expectedRequest: `GET /?website HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:r+KluO4YIu/eW9blDlcsGiUFpxc=
Date: %s`,
		response: `HTTP/1.1 200
x-oss-request-id: 50519080C4689A033D00235F
Date: Thu, 13 Sep 2012 07:51:28 GMT
Connection: close
Content-Length: 270
Server: AliyunOSS

<?xml version="1.0" encoding="UTF-8"?>
<WebsiteConfiguration xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
<IndexDocument>
<Suffix>index.html</Suffix>
    </IndexDocument>
    <ErrorDocument>
        <Key>error.html</Key>
    </ErrorDocument>
</WebsiteConfiguration>`,
		expectedResponse: &WebsiteConfiguration{
			IndexDocument: IndexDocument{
				Suffix: "index.html",
			},
			ErrorDocument: ErrorDocument{
				Key: "error.html",
			},
		},
	},

	{
		name: "GetBucketReferer",
		request: func(a *API) (interface{}, error) {
			r, err := a.GetBucketReferer(testBucketName)
			return r, err
		},
		expectedRequest: `GET /?referer HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:sPo+aFW5S/IPjbEmONETKOjDE2c=
Date: %s`,
		response: `HTTP/1.1 200
x-oss-request-id: 50519080C4689A033D00235F
Date: Thu, 13 Sep 2012 07:51:28 GMT
Connection: close
Content-Length: 242
Server: AliyunOSS

<?xml version="1.0" encoding="UTF-8"?>
<RefererConfiguration>
<AllowEmptyReferer>true</AllowEmptyReferer>
<RefererList>
<Referer>http://www.aliyun.com</Referer>
<Referer>https://www.aliyun.com</Referer>
</RefererList>
</RefererConfiguration>`,
		expectedResponse: &RefererConfiguration{
			AllowEmptyReferer: true,
			Referer: []string{
				"http://www.aliyun.com",
				"https://www.aliyun.com",
			},
		},
	},

	{
		name: "GetBucketLifecycle",
		request: func(a *API) (interface{}, error) {
			r, err := a.GetBucketLifecycle(testBucketName)
			return r, err
		},
		expectedRequest: `GET /?lifecycle HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:7aTx+kQARBcKul23R01R1Rw6yvU=
Date: %s`,
		response: `HTTP/1.1 200
x-oss-request-id: 534B372974E88A4D89060099
Date: Mon, 14 Apr 2014 01:17:29 GMT
Connection: close
Content-Length: 255
Server: AliyunOSS

<?xml version="1.0" encoding="UTF-8"?>
<LifecycleConfiguration>
  <Rule>
    <ID>delete after one day</ID>
    <Prefix>logs/</Prefix>
    <Status>Enabled</Status>
    <Expiration>
      <Days>1</Days>
    </Expiration>
  </Rule>
</LifecycleConfiguration>`,
		expectedResponse: &LifecycleConfiguration{
			Rule: []LifecycleRule{
				{
					ID:     "delete after one day",
					Prefix: "logs/",
					Status: "Enabled",
					Expiration: Expiration{
						Days: 1,
					},
				},
			},
		},
	},

	{
		name: "DeleteBucket",
		request: func(a *API) (interface{}, error) {
			return nil, a.DeleteBucket(testBucketName)
		},
		expectedRequest: `DELETE / HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:FXJsV//QmBZqMczqQkR2Lm+QUKY=
Date: %s`,
		response:         "HTTP/1.1 200\n",
		expectedResponse: nil,
	},

	{
		name: "DeleteBucketLogging",
		request: func(a *API) (interface{}, error) {
			return nil, a.DeleteBucketLogging(testBucketName)
		},
		expectedRequest: `DELETE /?logging HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:a/VOgtY0UksF67OgcuQZYx2OdZs=
Date: %s`,
		response:         "HTTP/1.1 200\n",
		expectedResponse: nil,
	},

	{
		name: "DeleteBucketWebsite",
		request: func(a *API) (interface{}, error) {
			return nil, a.DeleteBucketWebsite(testBucketName)
		},
		expectedRequest: `DELETE /?website HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:t0fB2Mn0sevKbSXp/c5q2XTGK3M=
Date: %s`,
		response:         "HTTP/1.1 200\n",
		expectedResponse: nil,
	},

	{
		name: "DeleteBucketLifecycle",
		request: func(a *API) (interface{}, error) {
			return nil, a.DeleteBucketLifecycle(testBucketName)
		},
		expectedRequest: `DELETE /?lifecycle HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:OImj+a+NC4oGhFL3lsSRj8iY6XI=
Date: %s`,
		response: `HTTP/1.1 204 No Content
x-oss-request-id: 534B372F74E88A4D89060124
Date: Mon, 14 Apr 2014 01:17:35 GMT
Connection: close
Content-Length: 0
Server: AliyunOSS
`,
		expectedResponse: nil,
	},

	{
		name: "PutObject",
		request: func(a *API) (interface{}, error) {
			return nil, a.PutObject(testBucketName, testObjectName, strings.NewReader("wefpofjwefew"))
		},
		expectedRequest: `PUT /object/name HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 12
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:gbvg8Xcdy0qvDT2e7uUdtj6/VZE=
Content-Type: application/octet-stream
Date: %s

wefpofjwefew`,
		response:         "HTTP/1.1 200 OK\n",
		expectedResponse: nil,
	},

	{
		name: "CopyObject",
		request: func(a *API) (interface{}, error) {
			r, err := a.CopyObject("source-bucket", "source-object", "target-bucket", "target-object")
			return r, err
		},
		expectedRequest: `PUT /target-object HTTP/1.1
Host: target-bucket.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 0
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:83kiZ9FOf79+NONgJAXk00Hzk4g=
Date: %s
X-Oss-Copy-Source: /source-bucket/source-object`,
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
		name: "GetObject",
		request: func(a *API) (interface{}, error) {
			w := new(bytes.Buffer)
			if err := a.GetObject(testBucketName, testObjectName, w); err != nil {
				return nil, err
			}
			if expected, actual := "abcdef", string(w.Bytes()); actual != expected {
				return nil, fmt.Errorf(expectBut, expected, actual)
			}
			return nil, nil
		},
		expectedRequest: `GET /object/name HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:pZol4Z6em1QCAz53w4OatKsdi3w=
Date: %s`,
		response: `HTTP/1.1 200 OK
x-oss-request-id: 3a89276f-2e2d-7965-3ff9-51c875b99c41
x-oss-object-type: Normal
Date: Fri, 24 Feb 2012 06:38:30 GMT
Last-Modified: Fri, 24 Feb 2012 06:07:48 GMT
ETag: "5B3C1A2E053D763E1B002CC607C5A0FE "
Content-Type: text/plain
Content-Length: 6
Connection: close
Server: AliyunOSS

abcdef`,
		expectedResponse: nil,
	},

	{
		name: "AppendObject",
		request: func(a *API) (interface{}, error) {
			r, err := a.AppendObject(testBucketName, testObjectName, strings.NewReader("sfweruewpinbeewa"), 0)
			return r, err
		},
		expectedRequest: `POST /object/name?append&position=0 HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 16
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:PwNSrS1NZvpPH6pfzPiIvQWH0G8=
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
		name: "DeleteObject",
		request: func(a *API) (interface{}, error) {
			return nil, a.DeleteObject(testBucketName, testObjectName)
		},
		expectedRequest: `DELETE /object/name HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:V1ehjYUAX1v6/ZUCzNKbCLKXQWE=
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
		expectedRequest: `POST /?delete HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 133
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:TmXuZ8SpwHR8lvCewR/jnnoR1qA=
Content-Md5: JDhCJwY/U4gHz5o/Q4eCoQ==
Date: %s

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
		name: "HeadObject",
		request: func(a *API) (interface{}, error) {
			r, err := a.HeadObject(testBucketName, testObjectName)
			return r, err
		},
		expectedRequest: `HEAD /object/name HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:B29hiJ0Fu10nq+kyeb4+vM6Cwns=
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
		name: "PutObjectACL",
		request: func(a *API) (interface{}, error) {
			return nil, a.PutObjectACL(testBucketName, testObjectName, PublicReadACL)
		},
		expectedRequest: `PUT /object/name?acl HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 0
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:Z31vuwnUDv9rezqnsdANa3+Utfs=
Date: %s
X-Oss-Acl: public-read`,
		response:         "HTTP/1.1 200 OK\n",
		expectedResponse: nil,
	},

	{
		name: "GetObjectACL",
		request: func(a *API) (interface{}, error) {
			r, err := a.GetObjectACL(testBucketName, testObjectName)
			return r, err
		},
		expectedRequest: `GET /object/name?acl HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:EbdVS7t4lipCbC4PPjUbwgOIToo=
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
		name: "InitUpload",
		request: func(a *API) (interface{}, error) {
			r, err := a.InitUpload(testBucketName, testObjectName)
			return r, err
		},
		expectedRequest: `POST /object/name?uploads HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 0
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:tuMEuJhbZGSP1wsWYvOI2awfPFw=
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
    <Bucket>bucket-name</Bucket>
    <Key>multipart.data</Key>
    <UploadId>0004B9894A22E5B1888A1E29F8236E2D</UploadId>
</InitiateMultipartUploadResult>`,
		expectedResponse: &InitiateMultipartUploadResult{
			Bucket:   "bucket-name",
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
		expectedRequest: `PUT /object/name?partNumber=1&uploadId=0004B9895DBBB6EC98E36 HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 16
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:MUWh2GXHDWeyeURcpCDFJIQdPlM=
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
		name: "UploadPartCopy",
		request: func(a *API) (interface{}, error) {
			r, err := a.UploadPartCopy(testBucketName, testObjectName, "0004B9895DBBB6EC98E36", 1, "source-bucket", "source-object")
			return r, err
		},
		expectedRequest: `PUT /object/name?partNumber=1&uploadId=0004B9895DBBB6EC98E36 HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 0
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:lY62AzamlvOzTWITUm1mPjp7D5Q=
Date: %s
X-Oss-Copy-Source: /source-bucket/source-object`,
		response: `HTTP/1.1 200 OK
Server: AliyunOSS
Content-Length: 232
Connection: close
x-oss-request-id: 3e6aba62-1eae-d246-6118-8ff42cd0c21a
Date: Thu, 17 Jul 2014 06:27:54 GMT'

<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
    <LastModified>2014-07-17T06:27:54.000Z</LastModified>
    <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE"</ETag>
</CopyPartResult>`,
		expectedResponse: &CopyPartResult{
			LastModified: parseTime(time.RFC3339Nano, "2014-07-17T06:27:54.000Z"),
			ETag:         `"5B3C1A2E053D763E1B002CC607C5A0FE"`,
		},
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
		expectedRequest: `POST /object/name?uploadId=0004B9895DBBB6EC98E36 HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 135
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:uhgZBIVEZXwHj/kbhGgqDH/q1yk=
Content-Md5: 9ZCUaLPTyBu1a7wYPPi23w==
Content-Type: application/octet-stream
Date: %s

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
		name: "AbortUpload",
		request: func(a *API) (interface{}, error) {
			return nil, a.AbortUpload(testBucketName, testObjectName, "0004B9895DBBB6EC98E36")
		},
		expectedRequest: `DELETE /object/name?uploadId=0004B9895DBBB6EC98E36 HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:KbFCcxzGSPnjmGmVuC8RSv0Pi1I=
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
		expectedRequest: `GET /?uploads HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:YcD3EXdEiXLUSdVkmxyE2lK14g0=
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
		expectedRequest: `GET /object/name?uploadId=0004B9895DBBB6EC98E36 HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:/JcXW8XD/ECBSA4cY5sIUi58J58=
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
		name: "PutBucketCORS",
		request: func(a *API) (interface{}, error) {
			return nil, a.PutBucketCORS(testBucketName, &CORSConfiguration{
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
		expectedRequest: `PUT /?cors HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Content-Length: 510
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:EOcmfX5IgAYU6yrHEgMsSyjrxLk=
Content-Md5: b3amOl4tM2Xd6JIKiHk9TQ==
Date: %s

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
		name: "GetBucketCORS",
		request: func(a *API) (interface{}, error) {
			r, err := a.GetBucketCORS(testBucketName)
			return r, err
		},
		expectedRequest: `GET /?cors HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:PM+ElDiB+ORFf9OqEhTVn5je2K0=
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

	{
		name: "DeleteBucketCORS",
		request: func(a *API) (interface{}, error) {
			return nil, a.DeleteBucketCORS(testBucketName)
		},
		expectedRequest: `DELETE /?cors HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Authorization: OSS ayahghai0juiSie:8krD9dvjeiyw4fwlYhHCkQwIfOM=
Date: %s`,
		response: `HTTP/1.1 204 No Content
x-oss-request-id: 5051845BC4689A033D0022BC
Date: Fri, 24 Feb 2012 05:45:34 GMT
Connection: close
Content-Length: 0
Server: AliyunOSS
`,
		expectedResponse: nil,
	},

	{
		name: "OptionObject",
		request: func(a *API) (interface{}, error) {
			r, err := a.OptionObject(testBucketName, testObjectName,
				AccessControlRequestMethod("PUT"),
				AccessControlRequestHeaders("x-oss-test"),
				Origin("http://www.example.com"))
			return r, err
		},
		expectedRequest: `OPTIONS /object/name HTTP/1.1
Host: bucket-name.oss-cn-hangzhou.aliyuncs.com
User-Agent: %s
Accept-Encoding: identity
Access-Control-Request-Headers: x-oss-test
Access-Control-Request-Method: PUT
Authorization: OSS ayahghai0juiSie:y5a3p8IEGw6n2bY9jUG4mu2ywVI=
Date: %s
Origin: http://www.example.com`,
		response: `HTTP/1.1 200 OK
x-oss-request-id: 5051845BC4689A033D0022BC
Access-Control-Allow-Origin: http://www.example.com
Access-Control-Allow-Methods: PUT
Access-Control-Expose-Headers: x-oss-test
Connection: close
`,
		expectedResponse: Header{
			"X-Oss-Request-Id":              []string{"5051845BC4689A033D0022BC"},
			"Access-Control-Allow-Origin":   []string{"http://www.example.com"},
			"Access-Control-Allow-Methods":  []string{"PUT"},
			"Access-Control-Expose-Headers": []string{"x-oss-test"},
		},
	},
}

func TestAllOssAPIs(t *testing.T) {
	for i := range apiTestcases {
		testAPI(t, &apiTestcases[i])
	}
}
func testAPI(t *testing.T, testcase *testcase) {
	server, err := NewMockServer(testcase.response)
	if err != nil {
		t.Fatalf(testcaseErr, testcase.name, err)
	}
	defer server.Close()
	api := New(testEndpoint, testID, testSecret, HTTPClient(&http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				// Hijack the address to 127.0.0.1 for testing
				return net.Dial(network, "127.0.0.1:"+server.Port())
			},
		},
	}))
	api.now = testTime
	response, err := testcase.request(api)
	if v := reflect.ValueOf(response); !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil()) {
		response = nil
	}
	if !reflect.DeepEqual(err, testcase.expectedError) {
		t.Fatalf(testcaseExpectBut, testcase.name, testcase.expectedError, err)
	}
	expectedRequest := fmt.Sprintf(testcase.expectedRequest, userAgent, testTimeText)
	if server.Err != nil {
		t.Fatalf(testcaseErr, testcase.name, err)
	}
	if server.Request != expectedRequest {
		t.Fatalf(testcaseExpectBut, testcase.name, expectedRequest, server.Request)
	}
	if !reflect.DeepEqual(response, testcase.expectedResponse) {
		t.Fatalf(testcaseExpectBut, testcase.name, testcase.expectedResponse, response)
	}
}

func TestSecurityToken(t *testing.T) {
	expected := "sec-token"
	api := New(testEndpoint, testID, testSecret, SecurityToken(expected))
	req, err := api.newRequest("GET", testBucketName, testObjectName, nil)
	if err != nil {
		t.Fatal(err)
	}
	if actual := req.Header.Get("X-Oss-Security-Token"); actual != expected {
		t.Fatalf(expectBut, expected, actual)
	}
}

func TestURLScheme(t *testing.T) {
	api := New(testEndpoint, testID, testSecret, URLScheme("https"))
	req, err := api.newRequest("GET", testBucketName, testObjectName, nil)
	if err != nil {
		t.Fatal(err)
	}
	if expected, actual := fmt.Sprintf("https://%s.%s/%s", testBucketName, testEndpoint, testObjectName), req.URL.String(); actual != expected {
		t.Fatalf(expectBut, expected, actual)
	}
}

func TestInvalidBucketName(t *testing.T) {
	api := New(testEndpoint, testID, testSecret)
	for _, bucket := range []string{
		"aBcde",  // capital letter
		"abcde_", // invalid char
		"-abcde", // not begin with lower caese letter or number
		"ab",     // too short
		strings.Repeat("a", 64), // too long
	} {
		if _, err := api.GetBucket(bucket); err != ErrInvalidBucketName {
			t.Fatalf(testcaseExpectBut, bucket, ErrInvalidBucketName, err)
		}
	}
}

func TestInvalidObjectName(t *testing.T) {
	api := New(testEndpoint, testID, testSecret)
	for _, object := range []string{
		"/abc",                    // start with /
		`\abc`,                    // start with \
		"abc\rde",                 // contains \r
		"abc\nde",                 // contains \n
		strings.Repeat("a", 1024), // too long
	} {
		if err := api.GetObject(testBucketName, object, new(bytes.Buffer)); err != ErrInvalidObjectName {
			t.Fatalf(testcaseExpectBut, object, ErrInvalidObjectName, err)
		}
	}
}

func TestIPEndpoint(t *testing.T) {
	api := New("127.0.0.1", testID, testSecret)
	req, err := api.newRequest("GET", testBucketName, testObjectName, nil)
	if err != nil {
		t.Fatal(err)
	}
	if expected, actual := "http://127.0.0.1/bucket-name/object/name", req.URL.String(); actual != expected {
		t.Fatalf(expectBut, expected, actual)
	}
}

func TestConnectionFail(t *testing.T) {
	errMessage := "injected failure"
	api := New(testEndpoint, testID, testSecret, HTTPClient(&http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return nil, errors.New(errMessage)
			},
		},
	}))
	if _, err := api.GetBucket(testBucketName); !strings.Contains(err.Error(), errMessage) {
		t.Fatalf(expectBut, errMessage, err.Error())
	}
}

func TestFailedOption(t *testing.T) {
	injectedFailure := errors.New("injected failure")
	failedOption := func(req *http.Request) error {
		return injectedFailure
	}
	api := New(testEndpoint, testID, testSecret)
	if _, err := api.GetService(failedOption); err != injectedFailure {
		t.Fatalf(expectBut, injectedFailure, err)
	}
}

func TestWrongResultType(t *testing.T) {
	panicked := false
	func() {
		defer func() {
			if err := recover(); err != nil {
				panicked = true
			}
		}()
		type wrongResult struct{} // wrongResult does not implement responseParser
		api := New(testEndpoint, testID, testSecret)
		api.handleResponse(&http.Response{
			StatusCode: 200,
		}, &wrongResult{})
	}()
	if !panicked {
		t.Fatal("expect panicking")
	}
}
