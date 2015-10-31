package oss

/*
import (
	"reflect"
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	errResp := `<?xml version="1.0" ?>
<Error xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
    <Code>AccessDenied</Code>
    <Message>Query-string authentication requires the Signature, Expires and OSSAccessKeyId parameters</Message>
    <RequestId>1D842BC5425544BB</RequestId>
    <HostId>oss-cn-hangzhou.aliyuncs.com</HostId>
</Error>`
	errXML, err := parseError(strings.NewReader(errResp))
	if err != nil {
		t.Fatal(err)
	}
	if expected := (&Error{
		Code:      "AccessDenied",
		Message:   "Query-string authentication requires the Signature, Expires and OSSAccessKeyId parameters",
		RequestID: "1D842BC5425544BB",
		HostID:    "oss-cn-hangzhou.aliyuncs.com",
	}); !reflect.DeepEqual(errXML, expected) {
		t.Fatalf(expectBut, expected, errXML)
	}
	if expected := "AccessDenied: Query-string authentication requires the Signature, Expires and OSSAccessKeyId parameters (1D842BC5425544BB, oss-cn-hangzhou.aliyuncs.com)"; errXML.Error() != expected {
		t.Fatalf(expectBut, expected, errXML.Error())
	}
}
*/
