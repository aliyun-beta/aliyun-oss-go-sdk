Extending the SDK
-----------------

With the clean and orthogonal implementation of the OSS Go SDK, it is rather
straightforward to extend the API.

### Define the body type

if the request's body is an XML, it can be sent via the oss.XMLBody option. You
just have to define the struct type, like below:

```go
	type LifecycleConfiguration struct {
		Rule []LifecycleRule
	}
	type LifecycleRule struct {
		ID         string
		Prefix     string
		Status     string
		Expiration Expiration
	}
```

### Define the result type

The result type must satisfy oss.ResponseParser interface. For example:

```go
func (r *LifecycleConfiguration) Parse(resp *http.Response) error {
	return xml.NewDecoder(resp.Body).Decode(r)
}
```

### Then use method oss.API.Do to do the work for you

```go
func (a *API) PutBucketLifecycle(bucket string, lifecycle *LifecycleConfiguration) error {
	return a.Do("PUT", bucket, "?lifecycle", nil, XMLBody(lifecycle))
}

func (a *API) GetBucketLifecycle(bucket string) (res *LifecycleConfiguration, _ error) {
	return res, a.Do("GET", bucket, "?lifecycle", &res)
}
```
