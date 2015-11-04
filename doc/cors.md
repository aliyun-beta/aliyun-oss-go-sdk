Cross-Origin Resource Sharing(CORS)
-----------------------------------

OSS provides the ability to configure CORS rules on a bucket.

### Set CORS rules

```go
	err := api.PutBucketCORS("bucket-name", &CORSConfiguration{
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
	if err != nil {
		log.Fatal(err)
	}
```

### Get CORS rules

```go
	res, err := api.GetBucketCORS("bucket-name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", res)
```

### Delete CORS rules

```go
	err := api.DeleteBucketCORS("bucket-name")
	if err != nil {
		log.Fatal(err)
	}
```
