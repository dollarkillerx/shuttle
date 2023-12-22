package conf

type Server struct {
	Protocol string `toml:"protocol" json:"protocol"`
	Addr     string `toml:"address" json:"addr"`
}

type HTTP struct {
	Path string `toml:"path" default:"/" json:"path"`
}

type WS struct {
	Path string `toml:"path" default:"/" json:"path"`
}

type ClientConfig struct {
	AgentToken string `toml:"agent_token"  json:"agent_token"`
	Server     `toml:"agent" json:"server"`
	HTTP       `toml:"http" json:"http"`
	WS         `toml:"ws" json:"ws"`
	Pac        bool `json:"pac"`
}
