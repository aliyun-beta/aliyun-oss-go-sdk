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
