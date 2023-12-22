package client

import (
	"fmt"
	"net"
	"testing"
)

func TestChinaIP(t *testing.T) {
	fmt.Println(isChinaIP("221.237.126.20"))
	fmt.Println(isChinaIP("43.135.75.195 "))
	fmt.Println(isChinaIP("192.227.234.228 "))
	ip := net.ParseIP("www.baid.com")
	if ip == nil {
		fmt.Println("Invalid IP address")
		return
	}

	if ip.IsLoopback() {
		fmt.Println(ip, "is loopback")
	}

	fmt.Println(ip.IsPrivate())

	if ip.IsGlobalUnicast() && !ip.IsLinkLocalUnicast() && !ip.IsInterfaceLocalMulticast() && !ip.IsLinkLocalMulticast() && !ip.IsMulticast() {
		fmt.Println(ip, "is a valid public IPv4 address")
	} else if ip.IsPrivate() {
		fmt.Println(ip, "is private")
	} else {
		fmt.Println(ip, "is not private and not public")
	}
}
