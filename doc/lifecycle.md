Object Lifecycle Management
---------------------------

OSS provides the ability to configure lifecycle management rules on a bucket.

### Set lifecycle rules

```go
	err := api.PutBucketLifecycle("bucket-name", &LifecycleConfiguration{
		Rule: []LifecycleRule{
			{
				ID:     "delete obsoleted files",
				Prefix: "obsoleted/",
				Status: "Enabled",
				Expiration: Expiration{
					Days: 3,
				},
			},
			{
				ID:     "delete temporary files",
				Prefix: "temporary/",
				Status: "Enabled",
				Expiration: Expiration{
					Date: parseTimePtr(time.RFC3339Nano, "2022-10-12T00:00:00.001Z"),
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
```

### Get Lifecycle rules

```go
	res, err := api.GetBucketLifecycle("bucket-name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", res)
```

### Delete Lifecycle rules

```go
	err := api.DeleteBucketLifecycle("bucket-name")
	if err != nil {
		log.Fatal(err)
	}
```
