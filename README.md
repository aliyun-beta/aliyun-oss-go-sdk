
Differences with Python SDK
---------------------------

* HTTP header keys are automatically converted into canonical format, e.g.
  x-oss-acl becomes X-Oss-Acl
* HTTP header User-Agent
* Go HTTP client does not support 100-Continue (https://github.com/golang/go/issues/3665)
* Go would do automatic charset detection for text/plain

Author
------

[Hǎiliàng Wáng](https://github.com/h12w)

## License

licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0.html)
