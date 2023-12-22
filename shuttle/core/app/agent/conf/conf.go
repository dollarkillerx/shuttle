package conf

import (
	"bytes"
	"os"

	"github.com/BurntSushi/toml"
)

type HTTP struct {
	Path string `toml:"path" default:"/"`
}

type WS struct {
	Path     string `toml:"path" default:"/"`
	Compress bool   `toml:"compress"`
}

type TLS struct {
	Cert string `toml:"cert"`
	Key  string `toml:"key"`
}

type Auth struct {
	Name string `toml:"name"`
	Pwd  string `toml:"pwd"`
}

type AgentConfig struct {
	Protocol        string       `toml:"protocol"`
	Listen          string       `toml:"listen"`           // 监听端口
	ManagerRPCAddr  string       `toml:"manager_rpc_addr"` // manager地址
	HTTP            HTTP         `toml:"http"`
	WS              WS           `toml:"ws"`
	TLS             TLS          `toml:"tls"`
	InternetAddress string       `toml:"internet_address"` // 服务外网地址
	AgentID         string       `toml:"agent_id"`         // agent id
	RegisterIP      string       `toml:"register_ip"`      // 服务注册ip
	RemoteSocks5    RemoteSocks5 `toml:"remote_socks5"`    // 再代理socks5
}

type RemoteSocks5 struct {
	Addr     string `toml:"addr"`
	UserName string `toml:"username"`
	Password string `toml:"password"`
}

func (a *AgentConfig) ReadConf() {
	agentToml, err := os.ReadFile("configs/agent.toml")
	if err != nil {
		panic(err)
	}

	_, err = toml.Decode(string(agentToml), a)
	if err != nil {
		panic(err)
	}
}

func (a *AgentConfig) WriteConf() {
	buffer := bytes.NewBuffer([]byte(""))
	encoder := toml.NewEncoder(buffer)
	err := encoder.Encode(a)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("agent.toml", buffer.Bytes(), 00666)
	if err != nil {
		panic(err)
	}
}
