package oss

import (
	"encoding/xml"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestListAllMyBucketsResult(t *testing.T) {
	xmlText := `<ListAllMyBucketsResult>
  <Owner>
    <ID>ut_test_put_bucket</ID>
    <DisplayName>ut_test_put_bucket</DisplayName>
  </Owner>
  <Buckets>
    <Bucket>
      <Location>oss-cn-hangzhou-a</Location>
      <Name>xz02tphky6fjfiuc0</Name>
      <CreationDate>2014-05-15T11:18:32.001Z</CreationDate>
    </Bucket>
    <Bucket>
      <Location>oss-cn-hangzhou-a</Location>
      <Name>xz02tphky6fjfiuc1</Name>
      <CreationDate>2014-05-15T11:18:32.002Z</CreationDate>
    </Bucket>
  </Buckets>
</ListAllMyBucketsResult>`
	xmlObject := ListAllMyBucketsResult{
		Owner: Owner{
			ID:          "ut_test_put_bucket",
			DisplayName: "ut_test_put_bucket",
		},
		Buckets: []Bucket{
			{
				Location:     "oss-cn-hangzhou-a",
				Name:         "xz02tphky6fjfiuc0",
				CreationDate: parseTime(time.RFC3339Nano, "2014-05-15T11:18:32.001Z"),
			},
			{
				Location:     "oss-cn-hangzhou-a",
				Name:         "xz02tphky6fjfiuc1",
				CreationDate: parseTime(time.RFC3339Nano, "2014-05-15T11:18:32.002Z"),
			},
		},
	}
	var r ListAllMyBucketsResult
	if err := xml.Unmarshal([]byte(xmlText), &r); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r, xmlObject) {
		t.Fatalf(expectBut, xmlObject, r)
	}
	buf, err := xml.MarshalIndent(xmlObject, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if string(buf) != xmlText {
		t.Fatalf(expectBut, xmlText, string(buf))
	}
}

func parseTime(layout, value string) time.Time {
	t, _ := time.Parse(layout, value)
	return t
}

func parseTimePtr(layout, value string) *time.Time {
	t, _ := time.Parse(layout, value)
	return &t
}

func TestAppendPosition(t *testing.T) {
	var pos AppendPosition
	err := pos.Parse(&http.Response{
		Header: http.Header{
			"X-Oss-Next-Append-Position": []string{"not int value"},
		},
	})
	if err == nil {
		t.Fatalf(expectBut, "error", err)
	}
}
