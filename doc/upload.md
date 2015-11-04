Multipart Upload
----------------

OSS supports upload a file via multiple API calls: Multipart Upload.

### Init multipart upload

```go
	initResult, err := api.InitUpload("bucket-name", "object/name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", initResult)
	uploadID = initResult.UploadID // remember this uploadID
```

### Upload parts

```go
	var parts []oss.Part
	for partNumber := 0; partNumber < maxPartNumber; partNumber++ {
		partcontents := fmt.Sprintf("contents of part %d\n", partNumber)
		partSize := int64(len(partContents))
		partResult, err := UploadPart("bucket-name", "object/name", uploadID, partNumber, strings.NewReader(partContents), partSize)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%#v\n", partResult)

		// collect parts information
		parts = append(parts, oss.Part{
			PartNumber: partNumber,
			LastModified: partResult.LastModified,
			ETag: partResult.ETag,
			Size: partSize,
		})
	}
```

### Complete upload after all parts are uploaded

```go
	res, err := api.CompleteUpload("bucket-name", "object/name", uploadID, &oss.CompleteMultipartUpload{
		Parts: parts,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", res)
```

### Abort the multipart upload

```go
	err := api.AbortUpload("bucket-name", "object/name", uploadID)
	if err != nil {
		log.Fatal(err)
	}
```
