package oss

import (
	"golang.org/x/sys/unix"
)

func Uname() *Utsname {
	var u unix.Utsname
	unix.Uname(&u)
	return &Utsname{
		SysName:    uToS(u.Sysname[:]),
		NodeName:   uToS(u.Nodename[:]),
		Release:    uToS(u.Release[:]),
		Version:    uToS(u.Version[:]),
		Machine:    uToS(u.Machine[:]),
		DomainName: uToS(u.Domainname[:]),
	}
}

func uToS(u []int8) string {
	buf := make([]byte, 0, 20)
	for _, c := range u {
		if c == 0 {
			break
		}
		buf = append(buf, byte(c))
	}
	return string(buf)
}
