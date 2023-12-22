package client

import (
	"bufio"
	"crypto/tls"
	"errors"
	"io"
	"net"

	"github.com/txthinking/socks5"
	"golang.org/x/net/proxy"
	"google.dev/google/shuttle/core/limits"
	"google.dev/google/shuttle/core/net_lib"
	"google.dev/google/shuttle/core/socks"
	"google.dev/google/shuttle/utils/log"
)

// Client holds contexts of the client
type Client struct {
	Config         *Config
	TLSConfig      *tls.Config
	Rules          *Rules
	localProxyAddr string
	proxySocks     proxy.Dialer
	socks5Proxy    *socks5.Server
}

// NewClient creates a client
func NewClient(localProxyAddr string) *Client {
	return &Client{
		Config:         &Config{},
		localProxyAddr: localProxyAddr,
	}
}

type _pd struct{}

func (p *_pd) Dial(network, addr string) (c net.Conn, err error) {
	return net_lib.DialTCP(network, "", addr)
}

// Serve starts the agent
func (c *Client) Serve(listener *net.TCPListener) (closeChannel chan struct{}, err error) {

	// socks5 proxy start
	socks5Proxy, err := socks5.NewClassicServer2(c.localProxyAddr, 0, 60)
	c.socks5Proxy = socks5Proxy
	go func() {
		err := socks5Proxy.ListenAndServe2(c)
		if err != nil {
			panic(err)
		}
	}()
	// socks5 proxy end

	// http proxy => socks5 proxy
	dial, err := proxy.SOCKS5("tcp", c.localProxyAddr, &proxy.Auth{
		User:     "",
		Password: "",
	}, &_pd{})
	if err != nil {
		return nil, err
	}
	if err := limits.Raise(); err != nil {
		log.Error("when", "try to raise system limits", "warning", err.Error())
	}
	c.proxySocks = dial
	// http proxy => socks5 proxy end

	closeChannel = make(chan struct{})

	go func() {
	loop:
		for {
			select {
			case <-closeChannel:
				break loop
			default:
				conn, err := listener.Accept()
				if err != nil {
					log.Errorf("Acceptance failed: %v", err)
					continue
				}

				go func() {
					br := bufio.NewReader(conn)
					// 获取当前协议 处理器
					handler, err := probeProtocol(br)
					if err != nil {
						conn.Close()
						log.Errorf("Probe protocol failed: %v", err)
						return
					}

					handler(c, &bufferedConn{conn, br})
				}()
			}
		}
	}()

	return closeChannel, nil
}

func probeProtocol(br *bufio.Reader) (func(*Client, net.Conn), error) {
	b, err := br.Peek(1)
	if err != nil {
		return nil, err
	}

	switch b[0] {
	case socks.Version:
		return (*Client).socks5Handler, nil
	default:
		return (*Client).httpHandler, nil
	}
}

type bufferedConn struct {
	net.Conn
	br *bufio.Reader
}

func (c *bufferedConn) Read(b []byte) (int, error) {
	return c.br.Read(b)
}

var protocol2wrapper = map[string]func(*Client, net.Conn) net.Conn{
	"ws":   (*Client).wrapWS,
	"wss":  (*Client).wrapWSS,
	"grpc": (*Client).wrapGRPC, // TODO: 还未实现
}

func (c *Client) dialServer() (net.Conn, error) {
	wrapper, ok := protocol2wrapper[c.Config.ServerProtocol]
	if !ok {
		return nil, errors.New("Unknow protocol")
	}

	conn, err := net.Dial("tcp", c.Config.ServerAddr)
	if err != nil {
		return nil, err
	}
	conn = wrapper(c, conn)

	// handshake
	if err := socks.WriteMethods([]byte{socks.MethodNoAuth}, conn); err != nil {
		conn.Close()
		return nil, err
	}
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		conn.Close()
		return nil, err
	}
	if buf[0] != socks.Version || buf[1] != socks.MethodNoAuth {
		conn.Close()
		return nil, errors.New("Handshake failed")
	}

	return conn, nil
}

// Config is the client configuration
type Config struct {
	Addr       string
	AgentToken string

	ServerProtocol string
	ServerAddr     string
	HTTPPath       string
	WSPath         string
	Pac            bool
	SetProxyOK     bool // 配置了proxy
}
