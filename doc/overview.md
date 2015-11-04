Overview
--------

Aliyun OSS (Object Storage Service) Go SDK is a client SDK to access Aliyun OSS
API, implemented in the Go programming language.

Here is a simple example about the basic functionalities that OSS Go SDK provides.

```go
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func main() {
	var accessKeyID, accessKeySecret string
	flag.StringVar(&accessKeyID, "id", "", "access key ID")
	flag.StringVar(&accessKeySecret, "secret", "", "access key Secret")
	flag.Parse()
	if accessKeyID == "" || accessKeySecret == "" {
		fmt.Println("go run main.go -id <your id> -secret <your secret>")
		flag.PrintDefaults()
		return
	}

	endPoint := "oss-cn-hangzhou.aliyuncs.com"

	// Create an API object
	api := oss.New(endPoint, accessKeyID, accessKeySecret)

	// List all buckets
	bucketList, err := api.GetService()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OK accessing OSS")
	log.Println(bucketList)

	// Create a bucket
	if err := api.PutBucket("bucket-name", oss.PrivateACL); err != nil {
		log.Fatal(err)
	}
	log.Println("bucket created or existed")

	// Upload an object
	testStr := "your test string"
	if err := api.PutObject("bucket-name", "object/name", strings.NewReader(testStr)); err != nil {
		log.Fatal(err)
	}
	log.Println("object uploaded")

	// Get the object
	buf := new(bytes.Buffer)
	headers, err := api.GetObject("bucket-name", "object/name", buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("object downloaded")
	log.Println(headers)

	if buf.String() != testStr {
		log.Fatalf("expectd %s but got %s", testStr, buf.String())
	}
	log.Println("object contents match")
}
```
