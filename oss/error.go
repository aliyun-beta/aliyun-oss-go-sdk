package oss

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidBucketName = errors.New("invalid bucket name")
	ErrInvalidObjectName = errors.New("invalid object name")
)

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
