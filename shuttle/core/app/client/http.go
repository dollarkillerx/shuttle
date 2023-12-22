package client

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/pkg/errors"
	"google.dev/google/shuttle/core/net_lib"
	"google.dev/google/shuttle/utils"
	"google.dev/google/shuttle/utils/log"
)

func (c *Client) wrapHTTPS(conn net.Conn) net.Conn {
	return c.wrapHTTP(tls.Client(conn, c.TLSConfig))
}

func (c *Client) wrapHTTP(conn net.Conn) net.Conn {
	return newHTTPWrapper(conn, c)
}

func isValidHTTPProxyRequest(req *http.Request) bool {
	if req.URL.Host == "" {
		return false
	}
	if req.Method != http.MethodConnect && req.URL.Scheme != "http" {
		return false
	}
	return true
}

func httpReply(statusCode int, status string) *http.Response {
	return &http.Response{
		ProtoMajor: 1,
		ProtoMinor: 1,
		StatusCode: statusCode,
		Status:     status,
	}
}

func (c *Client) httpHandler(conn net.Conn) {
	defer conn.Close()
	b := make([]byte, 0, 1024)
	for {
		var b1 [1024]byte
		n, err := conn.Read(b1[:])
		if err != nil {
			log.Error(err)
			return
		}
		b = append(b, b1[:n]...)
		if bytes.Contains(b, []byte{0x0d, 0x0a, 0x0d, 0x0a}) {
			break
		}
		if len(b) >= 2083+18 {
			log.Error(errors.New("HTTP header too long"))
			return
		}
	}

	bb := bytes.SplitN(b, []byte(" "), 3)
	if len(bb) != 3 {
		log.Error(errors.New("Invalid Request"))
		return
	}
	method, address := string(bb[0]), string(bb[1])
	var addr string
	if method == "CONNECT" {
		addr = address
	}
	if method != "CONNECT" {
		var err error
		addr, err = net_lib.GetAddressFromURL(address)
		if err != nil {
			log.Error(err)
			return
		}
	}

	tmp, err := c.proxySocks.Dial("tcp", addr)
	if err != nil {
		log.Error()
		return
	}
	rc := tmp.(*net.TCPConn)
	defer rc.Close()

	if method == "CONNECT" {
		_, err := conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			log.Error(err)
			return
		}
	}
	if method != "CONNECT" {
		if _, err := rc.Write(b); err != nil {
			log.Error(err)
			return
		}
	}
	go func() {
		var bf [1024 * 2]byte
		for {
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
		i, err := conn.Read(bf[:])
		if err != nil {
			return
		}
		if _, err := rc.Write(bf[0:i]); err != nil {
			return
		}
	}
	return
}

type httpWrapper struct {
	net.Conn
	client     *Client
	body       io.ReadCloser
	sentHeader bool

	ioBuf *bufio.Reader
	auth  string
}

func newHTTPWrapper(conn net.Conn, client *Client) *httpWrapper {
	//var auth string
	//cfg := client.Config
	//if cfg.Username != "" && cfg.Password != "" {
	//	s := base64.StdEncoding.EncodeToString([]byte(cfg.Username + ":" + cfg.Password))
	//	auth = "Basic " + s
	//}
	return &httpWrapper{
		Conn:   conn,
		client: client,
		ioBuf:  bufio.NewReader(conn),
		auth:   client.Config.AgentToken,
	}
}

func (h *httpWrapper) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}

	if h.body == nil {
		res, err := http.ReadResponse(h.ioBuf, nil)
		if err != nil {
			return 0, err
		}
		if res.StatusCode != 200 {
			res.Body.Close()
			return 0, fmt.Errorf("Response status is not OK: %s", res.Status)
		}
		if !utils.StrInSlice("chunked", res.TransferEncoding) {
			res.Body.Close()
			return 0, fmt.Errorf("Response is not chunked")
		}
		h.body = res.Body
	}

	return h.body.Read(b)
}

func (h *httpWrapper) Write(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	buf := bytes.NewBuffer(nil)
	if !h.sentHeader {
		buf.WriteString("POST ")
		buf.WriteString(h.client.Config.HTTPPath)
		buf.WriteString(" HTTP/1.1\r\n")
		buf.WriteString("Host: ")
		host, _, _ := net.SplitHostPort(h.client.Config.ServerAddr)
		buf.WriteString(host)
		buf.WriteString("\r\n")
		if h.auth != "" {
			buf.WriteString("Authorization: ")
			buf.WriteString(h.auth)
			buf.WriteString("\r\n")
		}
		buf.WriteString("Transfer-Encoding: chunked\r\n")
		buf.WriteString("\r\n")
		h.sentHeader = true
	}

	buf.WriteString(fmt.Sprintf("%X\r\n", len(b)))
	buf.Write(b)
	buf.WriteString("\r\n")
	if _, err := buf.WriteTo(h.Conn); err != nil {
		return 0, err
	}
	return len(b), nil
}
