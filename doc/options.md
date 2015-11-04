Optional Headers and Parameters
-------------------------------

Some methods of oss.API object provide optional named parameters for sending
additional header or URL parameter. The last optional parameter can either be
omitted, or provided with one or more option arguments.

These options are actually functions satisfying type Option:

```go
	type Option func(*http.Request) error
```

For example, oss.Expires is defined like below:

```go
func Expires(value string) Option {
	return setHeader("Expires", value)
}
func setHeader(key, value string) Option {
	return func(req *http.Request) error {
		req.Header.Set(key, value)
		return nil
	}
}
```

With these options defined, we are able to modify the http.Request object before
the SDK client sends the request, and multiple options can be applied via the
variadic function arguments.

```go
	err := api.PutObject("bucket-name", "object/name", f,
		oss.Expires("Fri, 28 Feb 2012 05:38:42 GMT"), // sets "Expires" header
		oss.Meta("user", "baymax"), // sets "X-Oss-Meta-User" header with value "baymax"
		)
	if err != nil {
		t.Fatal(err)
	}
```
