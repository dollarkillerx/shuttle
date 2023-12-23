package main

import (
	"google.dev/google/shuttle/core/app/agent"
	"google.dev/google/shuttle/core/app/agent/conf"
	"log"
)

func main() {
	var conf = new(conf.AgentConfig)
	conf.ReadConf()

	config := agent.FromAgentConfig(conf)
	tlsConfig, err := agent.GetServerTLSConfig(conf.TLS.Cert, conf.TLS.Key)
	if err != nil {
		log.Panicf("Get TLS configuration failed: %v", err)
		return
	}
	app, err := agent.NewApp(config, tlsConfig)
	app.Config.Verify = app.Verify
	if err != nil {
		log.Panicf("Failed to init Application: %v", err)
		return
	}
	app.Run()
}
