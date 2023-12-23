package agent

import (
	"github.com/txthinking/socks5"
	"google.dev/google/shuttle/core/socks"
	"google.dev/google/shuttle/utils"
	"log"

	"net"
	"strings"
)

type SocksHandler struct {
	conf *TaskConfig
}

func newSocksHandler(conf *TaskConfig) *SocksHandler {
	return &SocksHandler{
		conf: conf,
	}
}

func (s *SocksHandler) Handle(conn net.Conn) {
	defer conn.Close()

	// select method
	methods, err := socks.ReadMethods(conn)
	if err != nil {
		log.Printf(`[socks5] read methods failed: %v`, err)
		return
	}
	method := socks.MethodNoAcceptable
	for _, m := range methods {
		if m == socks.MethodNoAuth {
			method = m
		}
	}
	if err := socks.WriteMethod(method, conn); err != nil || method == socks.MethodNoAcceptable {
		if err != nil {
			log.Printf(`[socks5] write method failed: %v`, err)
		} else {
			log.Printf(`[socks5] methods is not acceptable`)
		}
		return
	}

	// read command
	request, err := socks.ReadRequest(conn)
	if err != nil {
		log.Printf(`[socks5] read command failed: %v \n`, err)
		return
	}
	switch request.Cmd {
	case socks.CmdConnect:
		s.handleConnect(conn, request)
	case socks.CmdBind:
		s.handleBind(conn, request)
	case socks.CmdUDP:
		// unsupported, since the agent based on TCP. using CmdUDPOverTCP instad.
		log.Printf(`[socks5] unsupported command CmdUDP`)
		if err := socks.NewReply(socks.CmdUnsupported, nil).Write(conn); err != nil {
			log.Printf(`[socks5] write reply failed: %v \n`, err)
		}
		return
	case socks.CmdUDPOverTCP:
		s.handleUDPOverTCP(conn, request)
	}
}

func (s *SocksHandler) handleConnect(conn net.Conn, req *socks.Request) {
	log.Printf(`[socks5] "connect" connect %s for %s`, req.Addr, conn.RemoteAddr())

	var newConn net.Conn
	var err error

	if strings.TrimSpace(s.conf.mountSocks5.Address) != "" {
		s5, err := socks5.NewClient(s.conf.mountSocks5.Address, s.conf.mountSocks5.Username, s.conf.mountSocks5.Password, 0, 60)
		if err != nil {
			log.Printf(`[socks5] "connect" dial remote failed: %v \n`, err)
			if err := socks.NewReply(socks.HostUnreachable, nil).Write(conn); err != nil {
				log.Printf(`[socks5] "connect" write reply failed: %v \n`, err)
			}
			return
		}

		newConn, err = s5.Dial("tcp", req.Addr.String())
		if err != nil {
			log.Printf(`[socks5] "connect" dial remote failed: %v \n`, err)
			if err := socks.NewReply(socks.HostUnreachable, nil).Write(conn); err != nil {
				log.Printf(`[socks5] "connect" write reply failed: %v \n`, err)
			}
			return
		}
	} else if strings.TrimSpace(s.conf.agentConf.RemoteSocks5.Addr) != "" {
		s5, err := socks5.NewClient(s.conf.agentConf.RemoteSocks5.Addr, s.conf.agentConf.RemoteSocks5.UserName, s.conf.agentConf.RemoteSocks5.Password, 0, 60)
		if err != nil {
			log.Printf(`[socks5] "connect" dial remote failed: %v \n`, err)
			if err := socks.NewReply(socks.HostUnreachable, nil).Write(conn); err != nil {
				log.Printf(`[socks5] "connect" write reply failed: %v \n`, err)
			}
			return
		}

		newConn, err = s5.Dial("tcp", req.Addr.String())
		if err != nil {
			log.Printf(`[socks5] "connect" dial remote failed: %v \n`, err)
			if err := socks.NewReply(socks.HostUnreachable, nil).Write(conn); err != nil {
				log.Printf(`[socks5] "connect" write reply failed: %v \n`, err)
			}
			return
		}
	} else {
		newConn, err = net.Dial("tcp", req.Addr.String())
		if err != nil {
			log.Printf(`[socks5] "connect" dial remote failed: %v \n`, err)
			if err := socks.NewReply(socks.HostUnreachable, nil).Write(conn); err != nil {
				log.Printf(`[socks5] "connect" write reply failed: %v \n`, err)
			}
			return
		}
	}

	defer newConn.Close()

	if err := socks.NewReply(socks.Succeeded, nil).Write(conn); err != nil {
		log.Printf(`[socks5] "connect" write reply failed: %v \n`, err)
		return
	}

	log.Printf(`[socks5] "connect" tunnel established %s <-> %s`, conn.RemoteAddr(), req.Addr)
	if err := utils.Transport(conn, newConn); err != nil {
		log.Printf(`[socks5] "connect" transport failed: %v \n`, err)
	}
	log.Printf(`[socks5] "connect" tunnel disconnected %s >-< %s`, conn.RemoteAddr(), req.Addr)
}

func (s *SocksHandler) handleBind(conn net.Conn, req *socks.Request) {
	log.Printf(`[socks5] "bind" bind for %s`, conn.RemoteAddr())
	listener, err := net.ListenTCP("tcp", nil)
	if err != nil {
		log.Printf(`[socks5] "bind" bind failed on listen: %v \n`, err)
		if err := socks.NewReply(socks.Failure, nil).Write(conn); err != nil {
			log.Printf(`[socks5] "bind" write reply failed %v \n`, err)
		}
		return
	}

	// first response: send listen address
	addr, _ := socks.NewAddrFromAddr(listener.Addr(), conn.LocalAddr())
	if err := socks.NewReply(socks.Succeeded, addr).Write(conn); err != nil {
		listener.Close()
		log.Printf(`[socks5] "bind" write reply failed %v \n`, err)
		return
	}

	newConn, err := listener.AcceptTCP()
	listener.Close()
	if err != nil {
		log.Printf(`[socks5] "bind" bind failed on accept: %v \n`, err)
		if err := socks.NewReply(socks.Failure, nil).Write(conn); err != nil {
			log.Printf(`[socks5] "bind" write reply failed %v \n`, err)
		}
		return
	}
	defer newConn.Close()

	// second response: accepted address
	raddr, _ := socks.NewAddr(newConn.RemoteAddr().String())
	if err := socks.NewReply(socks.Succeeded, raddr).Write(conn); err != nil {
		log.Printf(`[socks5] "bind" write reply failed %v \n`, err)
		return
	}

	log.Printf(`[socks5] "bind" tunnel established %s <-> %s`, conn.RemoteAddr(), newConn.RemoteAddr())
	if err := utils.Transport(conn, newConn); err != nil {
		log.Printf(`[socks5] "bind" transport failed: %v \n`, err)
	}
	log.Printf(`[socks5] "bind" tunnel disconnected %s >-< %s`, conn.RemoteAddr(), newConn.RemoteAddr())
}

func (s *SocksHandler) handleUDPOverTCP(conn net.Conn, req *socks.Request) {
	log.Printf(`[socks5] "udp-over-tcp" associate UDP for %s`, conn.RemoteAddr())
	udp, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Printf(`[socks5] "udp-over-tcp" UDP associate failed on listen: %v \n`, err)
		if err := socks.NewReply(socks.Failure, nil).Write(conn); err != nil {
			log.Printf(`[socks5] "udp-over-tcp" write reply failed %v \n`, err)
		}
		return
	}
	defer udp.Close()

	addr, _ := socks.NewAddrFromAddr(udp.LocalAddr(), conn.LocalAddr())
	if err := socks.NewReply(socks.Succeeded, addr).Write(conn); err != nil {
		log.Printf(`[socks5] "udp-over-tcp" write reply failed %v \n`, err)
		return
	}

	log.Printf(`[socks5] "udp-over-tcp" tunnel established %s <-> (UDP)%s`, conn.RemoteAddr(), udp.LocalAddr())
	if err := tunnelUDP(conn, udp); err != nil {
		log.Printf(`[socks5] "udp-over-tcp" tunnel UDP failed: %v \n`, err)
	}
	log.Printf(`[socks5] "udp-over-tcp" tunnel disconnected %s >-< (UDP)%s`, conn.RemoteAddr(), udp.LocalAddr())
}

func tunnelUDP(conn net.Conn, udp net.PacketConn) error {
	errc := make(chan error, 2)

	go func() {
		b := utils.LPool.Get().([]byte)
		defer utils.LPool.Put(b)

		for {
			n, addr, err := udp.ReadFrom(b)
			if err != nil {
				errc <- err
				return
			}

			saddr, _ := socks.NewAddr(addr.String())
			dgram := socks.NewUDPDatagram(
				socks.NewUDPHeader(uint16(n), 0, saddr), b[:n])
			if err := dgram.Write(conn); err != nil {
				errc <- err
				return
			}
		}
	}()

	go func() {
		for {
			dgram, err := socks.ReadUDPDatagram(conn)
			if err != nil {
				errc <- err
				return
			}

			addr, err := net.ResolveUDPAddr("udp", dgram.Header.Addr.String())
			if err != nil {
				continue
			}
			if _, err := udp.WriteTo(dgram.Data, addr); err != nil {
				errc <- err
				return
			}
		}
	}()

	return <-errc
}
