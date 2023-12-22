package net_lib

import (
	"net"
	"net/url"
	"time"

	"github.com/txthinking/socks5"
)

func ErrorReply(r *socks5.Request, c *net.TCPConn, e error) error {
	var p *socks5.Reply
	if r.Atyp == socks5.ATYPIPv4 || r.Atyp == socks5.ATYPDomain {
		p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv4, net.IPv4zero, []byte{0x00, 0x00})
	} else {
		p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv6, net.IPv6zero, []byte{0x00, 0x00})
	}
	if _, err := p.WriteTo(c); err != nil {
		return err
	}
	return e
}

func GetAddressFromURL(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	if _, _, err := net.SplitHostPort(u.Host); err == nil {
		return u.Host, nil
	}
	return net.JoinHostPort(u.Host, "80"), nil
}

func Conn2Conn(c, rc net.Conn, bufsize, timeout int) {
	go func() {
		bf := make([]byte, bufsize)
		for {
			if timeout != 0 {
				if err := rc.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
					return
				}
			}
			i, err := rc.Read(bf)
			if err != nil {
				return
			}
			if _, err := c.Write(bf[0:i]); err != nil {
				return
			}
		}
	}()
	bf := make([]byte, bufsize)
	for {
		if timeout != 0 {
			if err := c.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
				return
			}
		}
		i, err := c.Read(bf)
		if err != nil {
			return
		}
		if _, err := rc.Write(bf[0:i]); err != nil {
			return
		}
	}
	return
}
