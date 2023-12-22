package agent

import (
	"crypto/tls"
	"sync"
	"time"

	"google.dev/google/shuttle/pkg"
	"google.dev/google/shuttle/utils/log"
)

type Application struct {
	srv       Server
	client    RpcClient
	Config    *Config
	TLSConfig *tls.Config
}

func NewApp(config *Config, TLSConfig *tls.Config) (app *Application, err error) {
	var srv Server
	switch config.Protocol {
	case "ws":
		srv, err = newWsServer(config)
	case "wss":
		srv, err = newWssServer(config, TLSConfig)
	case "rpc":
		srv, err = newgRpcServer(config)
	case "https":
		srv, err = newHttpsServer(config, TLSConfig)
	default:
		srv, err = newHttpServer(config)
	}

	if err != nil {
		return nil, err
	}

	client := newRpcClient(config.conf.ManagerRPCAddr, config)
	return &Application{Config: config, TLSConfig: TLSConfig, srv: srv, client: *client}, nil
}

// Run start listen for client requests and keep alive to manager
func (app *Application) Run() {
	var group sync.WaitGroup

	group.Add(1)
	go func() {
		defer group.Done()
		err := app.client.RunAgentToManagerRpc()
		if err != nil {
			log.Errorf("RpcClient error: %v", err)
		}
	}()

	group.Add(1)
	go func() {
		defer group.Done()
		if err := app.srv.Serve(); err != nil {
			log.Errorf("Agent Server failed: %v", err)
		}
	}()

	group.Wait()
}

func (app *Application) Verify(token string) bool {
	if token == "" {
		return false
	}

	if app.Config.GetJWTAESKey() == "" {
		return false
	}

	var verifyToken = new(pkg.VerifyToken)
	verifyToken.FromToken(token, app.Config.GetJWTAESKey())

	if verifyToken.NodeID != app.Config.conf.AgentID {
		return false
	}

	if verifyToken.Expiration < time.Now().Unix() {
		return false
	}

	return true
}
