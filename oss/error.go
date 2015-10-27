package oss

import (
	"encoding/xml"
	"fmt"
	"io"
)

type ErrorXML struct {
	Code      string `xml:"Code"`
	Message   string `xml:"Message"`
	RequestID string `xml:"RequestId"`
	HostID    string `xml:"HostId"`
}

func (e *ErrorXML) Error() string {
	return fmt.Sprintf("%s: %s (%s, %s)", e.Code, e.Message, e.RequestID, e.HostID)
}

func parseErrorXML(rd io.Reader) (*ErrorXML, error) {
	e := new(ErrorXML)
	return e, xml.NewDecoder(rd).Decode(&e)
}
