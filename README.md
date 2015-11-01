Aliyun OSS Go SDK
=================

Quick Start
-----------

```bash
go get -u github.com/aliyun/aliyun-oss-go-sdk/oss
go test -v github.com/aliyun/aliyun-oss-go-sdk/oss
```

Highlights
----------
* Complete set of Aliyun OSS API
* 100% test coverage
* Idiomatic and flexible programming interface
* Clean and orthogonal implementation
* Named options for setting headers & parameters
* Extremely easy to extend for new APIs

Differences with Python SDK
---------------------------

* HTTP header User-Agent, e.g. aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
* Go HTTP client does not support 100-Continue (will be supported after Go 1.6, see https://github.com/golang/go/issues/3665)
* MIME type of txt file has default charset utf-8
* Go handles HTTP better than Python
  - Go GET request does not have redundant "Content-Length: 0" header
  - Parameters will be omitted if the argument is not set
  - HTTP header keys are automatically converted into canonical format, e.g.
    x-oss-acl becomes X-Oss-Acl
  - Go always sends URL parameters and headers in canonical order

Unclear issues in Spec
----------------------

* Resource or ResourceType or both in Error XML?

Author
------

[Hǎiliàng Wáng](https://github.com/h12w)

## License

licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0.html)
