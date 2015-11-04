Aliyun OSS Go SDK
=================

Aliyun OSS (Object Storage Service) Go SDK is a client SDK to access Aliyun OSS
API, implemented in the Go programming language.

Installation
------------

```bash
go get -u github.com/aliyun/aliyun-oss-go-sdk/oss
go test -v -cover github.com/aliyun/aliyun-oss-go-sdk/oss
```

Highlights
----------
* Complete set of Aliyun OSS API
* Thouroughly tested
  - 100% test coverage
  - intuitive table driven tests
  - full test suite completes within 2 seconds
* Lint clean
  - golint
  - go fmt
  - goimports
  - go vet
  - race detector
* Idiomatic & elegant
  - response is returned as a parsed object
  - error is returned as a Go error
  - named options for setting headers & parameters
* Great extensibility
  - clean and orthogonal implementation
  - users can easily extend the SDK when a new API method is supported
* No third party dependencies

Documentation
-------------

* [Overview](doc/overview.md)
* [API Object](doc/api-object.md)
* [Bucket](doc/bucket.md)
* [Object](doc/object.md)
* [Optional Headers and Parameters](doc/options.md)
* [Multipart Upload](doc/upload.md)
* [Cross-Origin Resource Sharing (CORS)](doc/cors.md)
* [Object Lifecycle Management](doc/lifecycle.md)
* [Extending the SDK](doc/extend.md)

Differences with Python SDK
---------------------------

* HTTP header User-Agent, e.g. aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
* Go HTTP client does not support 100-Continue (will be supported after Go 1.6, see https://github.com/golang/go/issues/3665)
* HTTP header keys are automatically converted into canonical format, e.g.
    x-oss-acl becomes X-Oss-Acl
* Go GET request does not have redundant "Content-Length: 0" header
* Parameters will be omitted if the argument is not set
* Go always sends URL parameters and headers in canonical order

Author
------

[Hǎiliàng Wáng](https://github.com/h12w)

License
-------

licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0.html)

References
----------
* [OSS Document (Chinese)](https://docs.aliyun.com/#/pub/oss)
