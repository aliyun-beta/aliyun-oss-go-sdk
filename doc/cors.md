Cross-Origin Resource Sharing(CORS)
-----------------------------------

OSS provides the ability to configure CORS rules on a bucket.

### Set CORS rules

```go
	err := api.PutBucketCORS("bucket-name", &oss.CORSConfiguration{
		CORSRule: []oss.CORSRule{
			{
				AllowedOrigin: []string{"xxxx"},
				AllowedMethod: []string{"xxxx"},
				AllowedHeader: []string{"xxxx"},
				ExposeHeader:  []string{"xxxx"},
				MaxAgeSeconds: 10000,
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
