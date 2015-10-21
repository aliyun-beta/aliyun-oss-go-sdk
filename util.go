package oss

import (
	"fmt"
	"runtime"
)

type Utsname struct {
	SysName    string
	NodeName   string
	Release    string
	Version    string
	Machine    string
	DomainName string
}

var userAgent string

func init() {
	uname := Uname()
	userAgent = fmt.Sprintf("aliyun-sdk-go/%s (%s/%s/%s;%s)", Version, uname.SysName, uname.Release, uname.Machine, runtime.Version())
}
