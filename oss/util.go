package oss

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
)

var pythonSDKCompatibleMode = false // comaptible with Python SDK (for unit testing only)

var userAgent = func() string {
	sys := getSysInfo()
	return fmt.Sprintf("aliyun-sdk-go/%s (%s/%s/%s;%s)", Version, sys.name, sys.release, sys.machine, runtime.Version())
}()

type sysInfo struct {
	name    string
	release string
	machine string
}

func getSysInfo() sysInfo {
	name := runtime.GOOS
	release := ""
	machine := runtime.GOARCH
	uname, err := exec.LookPath("uname")
	if err != nil {
		return sysInfo{name: runtime.GOOS, machine: runtime.GOARCH}
	}
	if out, err := exec.Command(uname, "-s").CombinedOutput(); err == nil {
		name = string(bytes.TrimSpace(out))
	}
	if out, err := exec.Command(uname, "-r").CombinedOutput(); err == nil {
		release = string(bytes.TrimSpace(out))
	}
	if out, err := exec.Command(uname, "-m").CombinedOutput(); err == nil {
		machine = string(bytes.TrimSpace(out))
	}
	return sysInfo{name: name, release: release, machine: machine}
}
