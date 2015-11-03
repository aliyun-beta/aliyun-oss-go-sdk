The OSS API Object
------------------

The OSS API object (oss.API) is the entry point to all the methods for manipulating
buckets and the objects in the buckets.

Rather than returning http.Response object directly, methods of oss.API return
parsed objects either from response body in XML or response headers. If an error
XML is returned, it is parsed and returned as a normal Go error interface{}.

Named optional parameters are provided on some of the methods, including the
oss.New function.

### Create oss.API

An oss.API object is created with oss.New funtion.

```go
	api := oss.New(endPoint, accessKeyID, accessKeySecret)
```

### Use HTTPS protocol

```go
	api := oss.New(endPoint, accessKeyID, accessKeySecret, oss.URLScheme("https"))
```

### Set the security token for Aliyun STS access

```go
	api := oss.New(endPoint, accessKeyID, accessKeySecret, oss.SecurityToken("your security token"))
```

### Set the underlying http.Client object for tuning parameters and more

```go
	api := oss.New(endPoint, accessKeyID, accessKeySecret,
		oss.HTTPClient(&http.Client{
			Timeout: 10*time.Second,
		}))
```

### Multiple optional arguments can be specified at the same time

```go
	api := oss.New(endPoint, accessKeyID, accessKeySecret,
		oss.URLScheme("https"),
		oss.HTTPClient(&http.Client{
			Timeout: 10*time.Second,
		}))
```
