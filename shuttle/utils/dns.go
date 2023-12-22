package utils

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
)

func LookupIP(domain string, dnsServer string) (string, error) {
	// 构造DNS请求
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	// 创建DNS客户端连接到指定的DNS服务器
	client := new(dns.Client)
	client.Net = "udp"
	client.Timeout = time.Second * 3 // 设置超时时间

	// 向DNS服务器发送请求并等待响应
	resp, _, err := client.Exchange(msg, dnsServer)
	if err != nil {
		return "", fmt.Errorf("failed to resolve domain name: %s", err.Error())
	}

	// 解析响应结果
	if len(resp.Answer) == 0 {
		return "", fmt.Errorf("no results found")
	}

	for _, ans := range resp.Answer {
		if a, ok := ans.(*dns.A); ok {
			return a.A.String(), nil
		}
	}

	return "", fmt.Errorf("failed to retrieve IP address")
}
