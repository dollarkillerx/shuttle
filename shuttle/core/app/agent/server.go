package agent

import (
	"crypto/tls"
	"sync"

	"google.dev/google/shuttle/core/app/agent/conf"
	"google.dev/google/shuttle/pkg"
	"google.dev/google/shuttle/utils"
	"google.dev/google/shuttle/utils/log"
)

type Server interface {
	Serve() error
}

// Config is the agent configuration
type Config struct {
	Protocol   string
	Addr       string
	Verify     func(token string) bool
	HTTPPath   string
	WSPath     string
	WSCompress bool
	Listen     string

	JWTAESKey string
	conf      *conf.AgentConfig

	mu sync.Mutex
}

type TaskConfig struct {
	agentConf   *conf.AgentConfig
	mountSocks5 *pkg.MountSocks5
}

func (c *Config) SetJWTAESKey(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.JWTAESKey = key
}

func (c *Config) GetJWTAESKey() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.JWTAESKey
}

func FromAgentConfig(conf *conf.AgentConfig) *Config {
	return &Config{
		Protocol:   conf.Protocol,
		Addr:       conf.ManagerRPCAddr,
		HTTPPath:   conf.HTTP.Path,
		WSPath:     conf.WS.Path,
		WSCompress: conf.WS.Compress,
		Listen:     conf.Listen,
		conf:       conf,
	}
}

func GetServerTLSConfig(cert, key string) (*tls.Config, error) {
	var certificate tls.Certificate
	var err error
	if cert == "" || key == "" {
		log.Info("Generate default TLS key pair")
		var rawCert, rawKey []byte
		rawCert, rawKey, err = utils.GenKeyPair()
		if err != nil {
			return nil, err
		}

		certificate, err = tls.X509KeyPair(rawCert, rawKey)
	} else {
		certificate, err = tls.LoadX509KeyPair(cert, key)
	}

	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{certificate}}, nil
}
