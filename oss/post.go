package oss

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path"
	"strings"
)

// PostOption is the option type for configuration multipart.Writer
type PostOption func(*multipart.Writer) error

// PostObject posts an object to OSS in MIME multipart format
func (a *API) PostObject(bucket, object, filename, policy string, options ...PostOption) (res Header, _ error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	policy = base64.StdEncoding.EncodeToString([]byte(policy))
	options = append(options, []PostOption{
		postObjectName(object),
		postAccessKeyID(a.accessKeyID),
		postPolicy(policy),
		postSignature(hmacSHA1([]byte(policy), []byte(a.accessKeySecret))),
		postFile(filename),
	}...)
	for _, option := range options {
		if err := option(w); err != nil {
			return nil, err
		}
	}
	w.Close()
	return res, a.Do("POST", bucket, "", &res, []Option{ContentType(w.FormDataContentType()), HTTPBody(buf)}...)
}

func setMultipartBoundary(boundary string) PostOption {
	return func(w *multipart.Writer) error {
		return w.SetBoundary(boundary)
	}
}

func postFile(filename string) PostOption {
	return func(w *multipart.Writer) error {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapeQuotes(path.Base(filename))))
		writer, _ := w.CreatePart(h)
		return safeCopy(writer, file)
	}
}
func safeCopy(w io.Writer, r io.Reader) error {
	if w == nil || r == nil {
		return errors.New("fail to copy")
	}
	_, err := io.Copy(w, r)
	return err
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// postPolicy is a PostOption to set policy
func postPolicy(value string) PostOption {
	return setMultipartField("policy", value)
}

// PostCacheControl is a PostOption to set Cache-Control
func PostCacheControl(value string) PostOption {
	return setMultipartField("Cache-Control", value)
}

// PostContentType is a PostOption to set Content-Type
func PostContentType(value string) PostOption {
	return setMultipartField("Content-Type", value)
}

// PostContentDisposition is a PostOption to set Content-Disposition
func PostContentDisposition(value string) PostOption {
	return setMultipartField("Content-Disposition", value)
}

// PostContentEncoding is a PostOption to set Content-Encoding
func PostContentEncoding(value string) PostOption {
	return setMultipartField("Content-Encoding", value)
}

// PostExpires is a PostOption to set Expires
func PostExpires(value string) PostOption {
	return setMultipartField("Expires", value)
}

// PostSuccessActionRedirect is a PostOption to set success_action_redirect
func PostSuccessActionRedirect(value string) PostOption {
	return setMultipartField("success_action_redirect", value)
}

// PostSuccessActionStatus is a PostOption to set success_action_status
func PostSuccessActionStatus(value string) PostOption {
	return setMultipartField("success_action_status", value)
}

// PostMeta is a PostOption to set X-Oss-Meta-* headers
func PostMeta(key, value string) PostOption {
	return setMultipartField("x-oss-meta-"+key, value)
}

// PostServerSideEncryption is a PostOption to set x-oss-server-side-encryption
func PostServerSideEncryption(value string) PostOption {
	return setMultipartField("x-oss-server-side-encryption", value)
}

// PostObjectACL is a PostOption to set x-oss-object-acl
func PostObjectACL(value ACLType) PostOption {
	return setMultipartField("x-oss-object-acl", string(value))
}

func postSignature(value string) PostOption {
	return setMultipartField("Signature", value)
}
func postObjectName(value string) PostOption {
	return setMultipartField("key", value)
}
func postAccessKeyID(value string) PostOption {
	return setMultipartField("OSSAccessKeyId", value)
}
func setMultipartField(key, value string) PostOption {
	return func(w *multipart.Writer) error {
		return w.WriteField(key, value)
	}
}
