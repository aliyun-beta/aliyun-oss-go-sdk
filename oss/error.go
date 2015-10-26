package oss

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
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
	if err := xml.NewDecoder(rd).Decode(&e); err != nil {
		return nil, err
	}
	e.Code = strings.TrimSpace(e.Code)
	e.Message = strings.TrimSpace(e.Message)
	e.RequestID = strings.TrimSpace(e.RequestID)
	e.HostID = strings.TrimSpace(e.HostID)
	return e, nil
}
