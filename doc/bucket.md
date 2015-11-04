Bucket
------

### Bucket name limitation

Here is the Go regular expression that specifies a valid bucket name.

```go
	regexp.MustCompile(`\A[a-z0-9][a-z0-9\-]{2,62}\z`)
```

Go SDK will check it before doing the real request to the server.

### Create a bucket

```go
	err := api.PutBucket("bucket-name", oss.PrivateACL)
	if err != nil {
		log.Fatal(err)
	}
```

### List all bucket of a user

```go
	result, err := api.GetService()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", result)
```

### Set a bucket's ACL

```go
	err := api.PutBucketACL("bucket-name", oss.PrivateACL)
	if err != nil {
		log.Fatal(err)
	}
```

### Get a bucket's ACL

```go
	acl, err := api.GetBucketACL("bucket-name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", acl)
```

### Create a bucket in a specified location

```go
	err := api.PutBucket("bucket-name", oss.BucketLocation("oss-cn-beijing"))
	if err != nil {
		log.Fatal(err)
	}
```

### Get a bucket's location

```go
	loc, err := api.GetBucketLocation("bucket-name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", loc)
```

### Delete a bucket

```go
	err := api.DeleteBucket("bucket-name")
	if err != nil {
		log.Fatal(err)
	}
```
