package client

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"sync"

	"google.dev/google/shuttle/core/app/client/conf"
	"google.dev/google/shuttle/utils/log"
)

type agentApiServer struct {
	listener     *net.TCPListener
	conf         conf.ClientConfig
	client       *Client
	closeChannel chan struct{}
	mu           sync.Mutex
}

func RouterRegister() {
	laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8986")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		panic(err)
	}

	log.Infof("Client starts to listen socks5://%s", listener.Addr().String())
	log.Infof("Client starts to listen http://%s", listener.Addr().String())

	client := NewClient(laddr.String())
	client.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	var aps = agentApiServer{
		listener: listener,
		client:   client,
	}
	// run core
	_, err = client.Serve(aps.listener)
	if err != nil {
		panic(err)
	}

	// 获取服务状态
	http.HandleFunc("/status", aps.apiStatus)
	// 下发任务
	http.HandleFunc("/link", aps.apiLink)
	// 结束任务
	http.HandleFunc("/stop", aps.apiStop)
}

func (a *agentApiServer) apiStatus(writer http.ResponseWriter, request *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.listener == nil {
		a.JSON(writer, 200, map[string]interface{}{
			"state": "deactivated",
		})
		return
	}

	a.JSON(writer, 200, map[string]interface{}{
		"conf":  a.conf,
		"state": "activated",
	})
}

func (a *agentApiServer) apiLink(writer http.ResponseWriter, request *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()

	var conf conf.ClientConfig
	err := a.BindJSON(request, &conf)
	if err != nil {
		a.JSON(writer, 400, map[string]interface{}{
			"message": "400",
		})
		return
	}

	if conf.Server.Addr == "" {
		a.JSON(writer, 400, map[string]interface{}{
			"message": "400",
		})
		return
	}

	a.client.Config.AgentToken = conf.AgentToken
	a.client.Config.ServerProtocol = conf.Server.Protocol
	a.client.Config.ServerAddr = conf.Server.Addr
	a.client.Config.HTTPPath = conf.HTTP.Path
	a.client.Config.WSPath = conf.WS.Path
	a.client.Config.SetProxyOK = true
	a.client.Config.Pac = conf.Pac

	a.conf = conf

	a.JSON(writer, 200, map[string]interface{}{
		"message": "success",
	})
	return
}

func (a *agentApiServer) apiStop(writer http.ResponseWriter, request *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.listener == nil {
		a.JSON(writer, 200, map[string]interface{}{
			"message": "success",
		})
		return
	}

	a.client.Config.SetProxyOK = false

	a.JSON(writer, 200, map[string]interface{}{
		"message": "success",
	})
}

func (a *agentApiServer) BindJSON(request *http.Request, v any) error {
	defer request.Body.Close()

	all, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(all, &v)
}

func (a *agentApiServer) JSON(write http.ResponseWriter, code int, v any) error {
	marshal, err := json.Marshal(v)
	if err != nil {
		return err
	}
	write.WriteHeader(code)
	write.Header().Set("content-type", "application/json")
	write.Write(marshal)

	return nil
}
