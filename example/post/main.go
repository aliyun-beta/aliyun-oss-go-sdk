package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"

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

	// Create a bucket
	if err := api.PutBucket("bucket-name", oss.PrivateACL); err != nil {
		log.Fatal(err)
	}
	log.Println("bucket created or existed")

	// Post an object
	policy := `{ "expiration": "2020-12-01T12:00:00.000Z", "conditions": [{"success_action_status": "200"}]}`
	postResult, err := api.PostObject("bucket-name", "posted/object", "testdata/test", policy, oss.PostSuccessActionStatus("200"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("object uploaded")
	fmt.Println(postResult)

	// Get the object
	buf := new(bytes.Buffer)
	headers, err := api.GetObject("bucket-name", "posted/object", buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("object downloaded")
	log.Println(headers)

	testStr := "sfweruewpinbeewa"
	if buf.String() != testStr {
		log.Fatalf("expectd %s but got %s", testStr, buf.String())
	}
	log.Println("object contents match")

	if err := api.DeleteObject("bucket-name", "posted/object"); err != nil {
		log.Fatal(err)
	}
	log.Println("object deleted")
}
