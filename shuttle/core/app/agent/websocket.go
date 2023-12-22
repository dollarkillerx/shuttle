package agent

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"google.dev/google/shuttle/pkg"
	"google.dev/google/shuttle/utils"
	"google.dev/google/shuttle/utils/log"
)

// WsServer accept request and call handler to handle it
type WsServer struct {
	Listener net.Listener
	conf     *Config
}

// newWsServer listen to address in conf, and return server instance
func newWsServer(conf *Config) (srv *WsServer, err error) {
	listener, err := net.Listen("tcp", conf.Listen)
	if err != nil {
		return nil, err
	}

	return &WsServer{
		Listener: listener,
		conf:     conf,
	}, nil
}

// todo: should return some kind if error
// Serve start to accept requests and call handler
func (w *WsServer) Serve() error {
	for {
		conn, err := w.Listener.Accept()
		if err != nil {
			log.Errorf("WsServer listener accept error: %v", err)
			continue
		}

		taskConfig := &TaskConfig{
			agentConf:   w.conf.conf,
			mountSocks5: &pkg.MountSocks5{},
		}
		handler := newWsHandler(conn, w.conf, taskConfig)
		go handler.handle(taskConfig)
	}

	return nil
}

// WsHandler handle the request and reply message
type WsHandler struct {
	net.Conn
	wsConn     *websocket.Conn
	conf       *Config
	buf        *bytes.Buffer
	ioBuf      *bufio.Reader
	upgrader   *websocket.Upgrader
	taskConfig *TaskConfig
}

// newWsHandler return a handler hold conf and conn
func newWsHandler(conn net.Conn, conf *Config, taskConfig *TaskConfig) *WsHandler {
	return &WsHandler{
		buf:   bytes.NewBuffer(make([]byte, 0, 1024)),
		ioBuf: bufio.NewReader(conn),
		upgrader: &websocket.Upgrader{
			EnableCompression: conf.WSCompress,
		},
		conf:       conf,
		Conn:       conn,
		taskConfig: taskConfig,
	}
}

// handle call SocksHandler
func (w *WsHandler) handle(conf *TaskConfig) {
	newSocksHandler(conf).Handle(w)
}

// Read implement net.Conn interface method
func (w *WsHandler) Read(b []byte) (n int, err error) {
	if w.wsConn == nil {
		w.wsConn, err = w.handshake()
		if err != nil {
			return
		}
	}

	if len(b) == 0 {
		return 0, nil
	}

	if w.buf.Len() > 0 {
		return w.buf.Read(b)
	}

	_, p, err := w.wsConn.ReadMessage()
	if err != nil {
		return 0, err
	}
	n = copy(b, p)
	w.buf.Write(p[n:])

	return
}

// Write implement net.Conn interface method
func (w *WsHandler) Write(b []byte) (n int, err error) {
	if w.wsConn == nil {
		w.wsConn, err = w.handshake()
		if err != nil {
			return
		}
	}

	if len(b) == 0 {
		return 0, nil
	}

	err = w.wsConn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

// handshake establish ws connection
func (w *WsHandler) handshake() (wconn *websocket.Conn, err error) {
	var req *http.Request
	for {
		req, err = http.ReadRequest(w.ioBuf)
		if err != nil {
			return
		}

		if w.conf.Verify != nil {
			if !utils.HttpBasicAuth(req.Header.Get("Authorization"), w.conf.Verify) {
				req.Body.Close()

				http4XXResponse(401).Write(w.Conn)
				continue
			}

			var verifyToken = new(pkg.VerifyToken)
			verifyToken.FromToken(req.Header.Get("Authorization"), w.conf.GetJWTAESKey())
			w.taskConfig.mountSocks5.Address = verifyToken.MountSocks5.Address
			w.taskConfig.mountSocks5.Username = verifyToken.MountSocks5.Username
			w.taskConfig.mountSocks5.Password = verifyToken.MountSocks5.Password
		}
		if !utils.StrEQ(req.URL.Path, w.conf.WSPath) ||
			req.Header.Get("Connection") != "Upgrade" ||
			req.Header.Get("Upgrade") != "websocket" {
			req.Body.Close()
			http4XXResponse(404).Write(w.Conn)
			continue
		}

		break
	}
	defer req.Body.Close()

	log.Infof("[websocket] upgrade request received: %s %s", req.Method, req.URL.Path)

	res := newHTTPRes4WS(w.Conn, bufio.NewReadWriter(w.ioBuf, bufio.NewWriter(w.Conn)))
	wconn, err = w.upgrader.Upgrade(res, req, nil)
	if err == nil {
		log.Info("[websocket] connection established")
	}
	return
}

// WssServer accept request and call handler to handle it
type WssServer struct {
	WsServer
	tlsConf *tls.Config
}

// newWsServer listen to address in conf, and return server instance
func newWssServer(conf *Config, tlsConf *tls.Config) (srv *WssServer, err error) {
	listener, err := net.Listen("tcp", conf.Listen)
	if err != nil {
		return nil, err
	}

	log.Infof("WSS %s", conf.Listen)
	return &WssServer{
		WsServer: WsServer{Listener: listener, conf: conf},
		tlsConf:  tlsConf,
	}, nil
}

// todo: should return some kind if error
// Serve start to accept requests and call handler
func (w *WssServer) Serve() error {
	for {
		conn, err := w.Listener.Accept()
		if err != nil {
			log.Errorf("WssServer listener accept error: %v", err)
			continue
		}

		taskConfig := &TaskConfig{
			agentConf:   w.conf.conf,
			mountSocks5: &pkg.MountSocks5{},
		}
		handler := newWssHandler(tls.Server(conn, w.tlsConf), w.conf, w.tlsConf, taskConfig)

		go handler.handle(taskConfig)
	}

	return nil
}

// WssHandler handle the request and reply message
// embed WsHandler to implement net.Conn interface
type WssHandler struct {
	WsHandler
	tlsConf *tls.Config
}

// newWssHandler return a handler hold conf, tls conf and conn
func newWssHandler(conn net.Conn, conf *Config, tlsConf *tls.Config, taskConfig *TaskConfig) *WssHandler {
	return &WssHandler{WsHandler: *newWsHandler(conn, conf, taskConfig), tlsConf: tlsConf}
}

// handle call SocksHandler
func (w *WssHandler) handle(conf *TaskConfig) {
	newSocksHandler(conf).Handle(w)
}

type httpRes4WS struct {
	proto         string
	header        http.Header
	contentLength int64
	statusCode    int

	wroteHeader bool
	written     int64
	conn        net.Conn
	ioBuf       *bufio.ReadWriter
}

func newHTTPRes4WS(conn net.Conn, ioBuf *bufio.ReadWriter) *httpRes4WS {
	r := &httpRes4WS{
		proto:         "HTTP/1.1",
		header:        http.Header{},
		contentLength: -1,

		wroteHeader: false,
		conn:        conn,
		ioBuf:       ioBuf,
	}
	return r
}

func (r *httpRes4WS) Header() http.Header {
	return r.header
}

func (r *httpRes4WS) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}

	r.written += int64(len(b))
	if r.contentLength != -1 && r.written > r.contentLength {
		return 0, http.ErrContentLength
	}

	return r.conn.Write(b)
}

func (r *httpRes4WS) WriteHeader(statusCode int) {
	if r.wroteHeader {
		return
	}
	r.wroteHeader = true

	r.statusCode = statusCode
	if cl := r.header.Get("Content-Length"); cl != "" {
		v, err := strconv.ParseInt(cl, 10, 64)
		if err == nil && v >= 0 {
			r.contentLength = v
		} else {
			r.header.Del("Content-Length")
		}
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	buf.WriteString(r.proto)
	buf.WriteString(" ")
	buf.WriteString(strconv.FormatInt(int64(r.statusCode), 10))
	buf.WriteString(" ")
	buf.WriteString(http.StatusText(r.statusCode))
	buf.WriteString("\r\n")
	r.header.Write(buf)
	buf.WriteString("\r\n")

	buf.WriteTo(r.conn)
}

func (r *httpRes4WS) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.conn, r.ioBuf, nil
}
