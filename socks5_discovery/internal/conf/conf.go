package conf

import (
	"github.com/BurntSushi/toml"

	"os"
)

type S5DiscoveryConfig struct {
	Listen string `toml:"listen"` // 监听端口
}

func (s *S5DiscoveryConfig) ReadConf() {
	s5DiscoveryToml, err := os.ReadFile("configs/s5_discovery.toml")
	if err != nil {
		panic(err)
	}

	_, err = toml.Decode(string(s5DiscoveryToml), s)
	if err != nil {
		panic(err)
	}
}
