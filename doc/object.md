Object
------

### Naming limitation

Here is the Go regular expression that specifies a valid object name. It should
be a valid UTF8 string and the length should be less than 1024.

```go
	regexp.MustCompile(`\A[^/\\](|[^\r\n]*)\z`)
```

Go SDK will check it before doing the real request to the server.

### Upload an object

PutObject reads from an io.Reader and upload the byte stream to the server. The
uploaded contents should be less than 5GB, and there are methods for multipart
uploads for file larger than 5G.

Upload an object from a file:

```go
	f, err := os.Open("your file name")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err := api.PutObject("bucket-name", "object/name", f)
	if err != nil {
		t.Fatal(err)
	}
```

Upload an object from a string:

```go
	err := api.PutObject("bucket-name", "object/name", strings.NewReader("your string"))
	if err != nil {
		t.Fatal(err)
	}
```

### Simulate a folder with an empty object ending with "/"

```go
	err := api.putobject("bucket-name", "object/as/a/directory/", bytes.newreader(nil))
	if err != nil {
		t.fatal(err)
	}
```

### Setting the HTTP headers of an object

One or more optional headers can be provided when necessary.

```go
	err := api.PutObject("bucket-name", "object/name", f,
		oss.Expires("Fri, 28 Feb 2012 05:38:42 GMT"), // sets "Expires" header
		oss.Meta("user", "baymax"), // sets "X-Oss-Meta-User" header with value "baymax"
		)
	if err != nil {
		t.Fatal(err)
	}
```

### Append to an object

```go
	appendPos, err := api.AppendObject("bucket-name", "object/name",
		strings.NewReader("appended contents"),
		oss.AppendPosition(123))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(appendPos)
```

### List all objects in a bucket

```go
	// list all objects
	res, err := api.GetBucket("bucket-name")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", res)
```

### List objects in a bucket satisfying some options

```go
	// list all objects
	res, err := api.GetBucket("bucket-name", oss.Prefix("pic"), oss.Delimiter("/"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", res)
```

### Get an object

Get the contents of an object and its associated header.

```go
	buf := new(bytes.Buffer)
	headers, err := api.GetObject("bucket-name", "object/name", buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf.String())
	fmt.Println(headers)
```

### Get an object with options

For example, read bytes 20 to 100 from an object.

```go
	buf := new(bytes.Buffer)
	headers, err := api.GetObject("bucket-name", "object/name", buf,
		oss.Range("bytes=20-100"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf.String())
```

### Save an object to a file

```go
	f, err := os.Create("file_name")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err := api.GetObject("bucket-name", "object/name", f)
	if err != nil {
		log.Fatal(err)
	}
```

### Get object headers only

```go
	headers, err := api.HeadObject("bucket-name", "object/name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(headers)
```

### Delete an Object

```go
	if err := api.DeleteObject("bucket-name", "object/name"); err != nil {
		log.Fatal(err)
	}
```

### Delete a number of objects at the same time

```go
	if err := api.DeleteObjects("bucket-name", "object1", "object2", "object3"); err != nil {
		log.Fatal(err)
	}
```

### Copy an object online

```go
	res, err := api.CopyObject("source-bucket", "source/object", "target-bucket", "target/object")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", res)
```

### Modify an object's meta headers by copying to the same location

```go
	bucket, object := "bucket-name", "object/name"
	res, err := api.CopyObject(bucket, object, bucket, object, oss.ContentType("image/png"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", res)
```
