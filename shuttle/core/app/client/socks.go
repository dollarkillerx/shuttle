package client

import (
	"github.com/txthinking/socks5"
	"google.dev/google/shuttle/core/socks"
	"google.dev/google/shuttle/utils"
	"log"

	"io"
	"net"
	"time"
)

func (c *Client) wrapSocks(conn net.Conn) net.Conn {
	return conn
}

func (c *Client) socks5Handler(conn net.Conn) {
	defer conn.Close()

	if err := c.socks5Proxy.Negotiate(conn); err != nil {
		log.Println(err)
		return
	}
	r, err := c.socks5Proxy.GetRequest(conn)
	if err != nil {
		log.Println(err)
		return
	}
	if err := c.socks5Proxy.Handle.TCPHandle(c.socks5Proxy, conn, r); err != nil {
		log.Println(err)
	}
}

func (c *Client) TCPHandle(server *socks5.Server, conn net.Conn, request *socks5.Request) error {
	if request.Cmd == socks5.CmdConnect {
		c.handleConnect(server, conn, request)

		rc, err := request.Connect(conn)
		if err != nil {
			return err
		}
		defer rc.Close()
		go func() {
			var bf [1024 * 2]byte
			for {
				if server.TCPTimeout != 0 {
					if err := rc.SetDeadline(time.Now().Add(time.Duration(server.TCPTimeout) * time.Second)); err != nil {
						return
					}
				}
				i, err := rc.Read(bf[:])
				if err != nil {
					return
				}
				if _, err := conn.Write(bf[0:i]); err != nil {
					return
				}
			}
		}()
		var bf [1024 * 2]byte
		for {
			if server.TCPTimeout != 0 {
				if err := conn.SetDeadline(time.Now().Add(time.Duration(server.TCPTimeout) * time.Second)); err != nil {
					return nil
				}
			}
			i, err := conn.Read(bf[:])
			if err != nil {
				return nil
			}
			if _, err := rc.Write(bf[0:i]); err != nil {
				return nil
			}
		}
		return nil
	}
	if request.Cmd == socks5.CmdUDP {
		caddr, err := request.UDP(conn, server.ServerAddr)
		if err != nil {
			return err
		}
		ch := make(chan byte)
		defer close(ch)
		server.AssociatedUDP.Set(caddr.String(), ch, -1)
		defer server.AssociatedUDP.Delete(caddr.String())
		io.Copy(io.Discard, conn)
		return nil
	}
	return socks5.ErrUnsupportCmd
}

func (c *Client) UDPHandle(server *socks5.Server, addr *net.UDPAddr, datagram *socks5.Datagram) error {
	handle := new(socks5.DefaultHandle)
	return handle.UDPHandle(server, addr, datagram)
}

func (c *Client) handleConnect(server *socks5.Server, conn net.Conn, req *socks5.Request) {
	var nextHop net.Conn
	var err error
	var isProxy bool

	if PacProxy(req.Host(), c.Config.Pac) && c.Config.SetProxyOK == true && c.Config.ServerAddr != "" {
		log.Printf(`[socks5] "connect" dial agent to connect %s for %s \n`, req.Address(), conn.RemoteAddr())

		isProxy = true
		nextHop, err = c.dialServer()
		if err != nil {
			log.Printf(`[socks5] "connect" dial agent failed: %v \n`, err)
			if err = socks.NewReply(socks.HostUnreachable, nil).Write(conn); err != nil {
				log.Printf(`[socks5] "connect" write reply failed: %v \n`, err)
			}
			return
		}
		defer nextHop.Close()

	} else {
		log.Printf(`[socks5] "connect" dial %s for %s \n`, req.Address(), conn.RemoteAddr())

		nextHop, err = net.Dial("tcp", req.Address())
		if err != nil {
			log.Printf(`[socks5] "connect" dial remote failed: %v \n`, err)
			if err = socks.NewReply(socks.HostUnreachable, nil).Write(conn); err != nil {
				log.Printf(`[socks5] "connect" write reply failed: %v \n`, err)
			}
			return
		}
		defer nextHop.Close()
	}

	var dash rune
	if isProxy {
		if err = req.Write(nextHop); err != nil {
			log.Printf(`[socks5] "connect" send request failed: %v \n`, err)
			return
		}
		dash = '-'
	} else {
		if err = socks.NewReply(socks.Succeeded, nil).Write(conn); err != nil {
			log.Printf(`[socks5] "connect" write reply failed: %v \n`, err)
			return
		}
		dash = '='
	}

	log.Printf(`[socks5] "connect" tunnel established %s <%c> %s \n`, conn.RemoteAddr(), dash, req.Address())
	if err := utils.Transport(conn, nextHop); err != nil {
		log.Printf(`[socks5] "connect" transport failed: %v \n`, err)
	}
	log.Printf(`[socks5] "connect" tunnel disconnected %s >%c< %s \n`, conn.RemoteAddr(), dash, req.Address())
}
