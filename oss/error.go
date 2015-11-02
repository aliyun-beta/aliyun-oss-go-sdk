package oss

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrInvalidBucketName happens when the bucket name is not valid
	ErrInvalidBucketName = errors.New("invalid bucket name")
	// ErrInvalidObjectName happens when the object name is not valid
	ErrInvalidObjectName = errors.New("invalid object name")
)

// Error represents the XML error returned by OSS APIs
type Error struct {
	Code         string
	Message      string
	RequestID    string `xml:"RequestId"`
	Resource     string
	ResourceType string
	HostID       string `xml:"HostId"`

	HTTPStatusCode int    `xml:"-"`
	HTTPStatus     string `xml:"-"`
	ParseError     error  `xml:"-"`
}

func (e *Error) Error() string {
	if e.ParseError != nil {
		return fmt.Sprintf("%s: (%s)", e.HTTPStatus, e.ParseError.Error())
	}
	return fmt.Sprintf("%s (%s): %s (%s, %s)", e.Code, e.HTTPStatus, e.Message, e.RequestID, e.HostID)
}

func parseError(resp *http.Response) error {
	errObj := new(Error)
	errObj.ParseError = xml.NewDecoder(resp.Body).Decode(&errObj)
	errObj.HTTPStatusCode, errObj.HTTPStatus = resp.StatusCode, resp.Status
	return errObj
}
