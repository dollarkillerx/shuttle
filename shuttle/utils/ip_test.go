package utils

import "testing"

func TestGetIP(t *testing.T) {
	ips := GetIP()
	t.Log(ips)
}

func TestGetMacAddrs(t *testing.T) {
	macs := GetMacAddrs()
	t.Log(macs)
}
