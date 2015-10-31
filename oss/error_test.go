package oss

import (
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	errResp := &http.Response{
		StatusCode: 501,
		Status:     "501 Not Implemented",
		Body: ioutil.NopCloser(strings.NewReader(`<?xml version="1.0" ?>
<Error xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
    <Code>NotImplemented</Code>
    <Message>A header you provided implies functionality that is not implemented.</Message>
    <Header>If-Modified-Since</Header>
    <RequestId>77E534EBF90372BE</RequestId>
    <HostId>oss-cn-hangzhou.aliyuncs.com</HostId>
</Error>`)),
	}

	errObj := parseError(errResp)
	if expected := (&Error{
		Code:           "NotImplemented",
		Message:        "A header you provided implies functionality that is not implemented.",
		RequestID:      "77E534EBF90372BE",
		HostID:         "oss-cn-hangzhou.aliyuncs.com",
		HTTPStatusCode: 501,
		HTTPStatus:     "501 Not Implemented",
	}); !reflect.DeepEqual(errObj, expected) {
		t.Fatalf(expectBut, expected, errObj)
	}
	if expected := "NotImplemented (501 Not Implemented): A header you provided implies functionality that is not implemented. (77E534EBF90372BE, oss-cn-hangzhou.aliyuncs.com)"; errObj.Error() != expected {
		t.Fatalf(expectBut, expected, errObj.Error())
	}
}

func TestUnparsableError(t *testing.T) {
	errObj := &Error{
		HTTPStatusCode: 501,
		HTTPStatus:     "501 Not Implemented",
		ParseError:     errors.New("XML syntax error on line 2: unquoted or missing attribute value in element"),
	}
	if expected := "501 Not Implemented: (XML syntax error on line 2: unquoted or missing attribute value in element)"; errObj.Error() != expected {
		t.Fatalf(expectBut, expected, errObj.Error())
	}
}
