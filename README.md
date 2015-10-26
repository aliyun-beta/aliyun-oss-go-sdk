
Differences with Python SDK
---------------------------

* HTTP header User-Agent, e.g. aliyun-sdk-go/0.1.1 (Linux/3.16.0-51-generic/x86_64;go1.5.1)
* Go HTTP client does not support 100-Continue (will be supported after Go 1.6, see https://github.com/golang/go/issues/3665)
* Go handles HTTP better than Python
  - Go GET request does not have Content-Length header
  - Go would do automatic charset detection for text/plain
  - Parameters will be omitted if the argument is not set
  - HTTP header keys are automatically converted into canonical format, e.g.
    x-oss-acl becomes X-Oss-Acl

Author
------

[Hǎiliàng Wáng](https://github.com/h12w)

## License

licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0.html)
