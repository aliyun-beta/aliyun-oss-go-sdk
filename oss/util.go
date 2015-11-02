package oss

import (
	"fmt"
	"runtime"
)

var (
	userAgent               string
	pythonSDKCompatibleMode = false // comaptible with Python SDK (for unit testing only)
)

func init() {
	uname := getUname()
	userAgent = fmt.Sprintf("aliyun-sdk-go/%s (%s/%s/%s;%s)", Version, uname.SysName, uname.Release, uname.Machine, runtime.Version())
}

type utsname struct {
	SysName    string
	NodeName   string
	Release    string
	Version    string
	Machine    string
	DomainName string
}
