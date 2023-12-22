package conf

import (
	"google.dev/google/common/pkg/conf"
	"google.dev/google/shuttle/utils"
)

// CONFIG  global configuration
var CONFIG Configuration

// Configuration ...
type Configuration struct {
	GraphQLListenAddress string
	RPCListenAddress     string
	AgentAesKey          string
	Debug                bool
	EnablePlayground     bool
	JWTConfiguration     JWTConfiguration
	TimeOut              int
	LoggerConfig         conf.LoggerConfig

	Socks5DiscoveryAddress string

	PostgresConfiguration conf.PostgresConfiguration
	ETCDConfiguration     conf.ETCDConfiguration
	Salt                  string

	CORSAllowedOrigins []string
	StaticServerConfig StaticServerConfiguration
}

// StaticServerConfiguration 静态资源服务配置
type StaticServerConfiguration struct {
	EntityTagsFileName string
	StaticPath         string
}

var defaultStaticServerConfig = StaticServerConfiguration{
	StaticPath:         "static",
	EntityTagsFileName: ".entity_tags.json",
}

// GetDefault ...
func (s *StaticServerConfiguration) GetDefault() *StaticServerConfiguration {
	return &defaultStaticServerConfig
}

// JWTConfiguration ...
type JWTConfiguration struct {
	SecretKey    string
	OperationKey string
}

func InitConfiguration(configName string, configPath string) error {
	var c Configuration
	err := conf.InitConfiguration(configName, []string{configPath}, &c)
	c.AgentAesKey = utils.PadKeyString(c.AgentAesKey, 32)
	CONFIG = c
	return err
}
