package agent

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"google.dev/google/shuttle/pkg"
	"io"
	"net"
	"net/http"

	"google.dev/google/shuttle/utils"
	"google.dev/google/shuttle/utils/log"
)

// HttpServer accept request and call handler to handle it
type HttpServer struct {
	conf     *Config
	Listener net.Listener
}

// newHttpServer listen to address in conf, and return server instance
func newHttpServer(conf *Config) (*HttpServer, error) {
	listener, err := net.Listen("tcp", conf.Listen)
	if err != nil {
		return nil, err
	}
	return &HttpServer{conf: conf, Listener: listener}, nil
}

// todo: should return some kind if error
// Serve start to accept requests and call handler
func (s *HttpServer) Serve() error {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Errorf("HttpServer listener accept error: %v", err)
			continue
		}

		taskConfig := TaskConfig{
			agentConf:   s.conf.conf,
			mountSocks5: &pkg.MountSocks5{},
		}
		handler := newHttpHandler(conn, s.conf)
		go handler.handle(&taskConfig)
	}

	return nil
}

// HttpsServer accept request and call handler to handle it
type HttpsServer struct {
	conf     *Config
	tlsConf  *tls.Config
	Listener net.Listener
}

// newHttpsServer listen to address in conf, and return server instance
func newHttpsServer(conf *Config, tlsConf *tls.Config) (*HttpsServer, error) {
	listener, err := net.Listen("tcp", conf.Listen)
	if err != nil {
		return nil, err
	}
	return &HttpsServer{conf: conf, Listener: listener, tlsConf: tlsConf}, nil
}

// todo: should return some kind if error
// Serve start to accept requests and call handler
func (s *HttpsServer) Serve() error {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Errorf("HttpsServer listener accept error: %v", err)
			continue
		}

		taskConfig := TaskConfig{
			agentConf:   s.conf.conf,
			mountSocks5: &pkg.MountSocks5{},
		}
		handler := newHttpsHandler(conn, s.conf, s.tlsConf)
		go handler.handle(&taskConfig)
	}

	return nil
}

// HttpHandler handle the request and reply message
type HttpHandler struct {
	net.Conn
	body       io.ReadCloser
	sentHeader bool
	ioBuf      *bufio.Reader
	conf       *Config
}

// newHttpHandler return a handler hold conf and conn
func newHttpHandler(conn net.Conn, conf *Config) *HttpHandler {
	return &HttpHandler{
		Conn:  conn,
		ioBuf: bufio.NewReader(conn),
		conf:  conf,
	}
}

// handle call SocksHandler
func (h *HttpHandler) handle(conf *TaskConfig) {
	newSocksHandler(conf).Handle(h)
}

// Read implement net.Conn interface
func (h *HttpHandler) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}

	for h.body == nil {
		req, err := http.ReadRequest(h.ioBuf)
		if err != nil {
			return 0, err
		}
		if h.conf.Verify != nil {
			if !utils.HttpBasicAuth(req.Header.Get("Authorization"), h.conf.Verify) {
				req.Body.Close()
				http4XXResponse(401).Write(h.Conn)
				continue
			}
		}
		if !utils.StrEQ(req.URL.Path, h.conf.HTTPPath) {
			req.Body.Close()
			http4XXResponse(404).Write(h.Conn)
			continue
		}
		if !utils.StrInSlice("chunked", req.TransferEncoding) {
			req.Body.Close()
			http4XXResponse(400).Write(h.Conn)
			continue
		}
		h.body = req.Body
	}

	return h.body.Read(b)
}

// Write implement net.Conn interface
func (h *HttpHandler) Write(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	buf := bytes.NewBuffer(nil)
	if !h.sentHeader {
		buf.WriteString("HTTP/1.1 200 OK\r\n")
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

// HttpsHandler handle the request and reply message
// embed HttpHandler to implement net.Conn interface
type HttpsHandler struct {
	HttpHandler
	tlsConf *tls.Config
}

// newHttpsHandler return a handler hold conf and conn
func newHttpsHandler(conn net.Conn, conf *Config, tlsConf *tls.Config) *HttpsHandler {
	return &HttpsHandler{
		HttpHandler: *newHttpHandler(conn, conf),
		tlsConf:     tlsConf,
	}
}

// handle call SocksHandler
func (h *HttpsHandler) handle(conf *TaskConfig) {
	newSocksHandler(conf).Handle(tls.Server(h.Conn, h.tlsConf))
}

func http4XXResponse(code int) *http.Response {
	body := bytes.NewBufferString(
		fmt.Sprintf("<h1>%d</h1><p>%s<p>", code, http.StatusText(code)))
	header := make(http.Header)
	if code == 401 {
		header.Add("WWW-Authenticate", `Basic realm="auth"`)
	}
	return &http.Response{
		StatusCode:    code,
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: int64(body.Len()),
		Body:          io.NopCloser(body),
		Header:        header,
	}
}
