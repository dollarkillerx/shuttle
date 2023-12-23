package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.dev/google/shuttle/proto/manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RpcClient connect to manager and keep alive
type RpcClient struct {
	target string
	Config *Config
}

func newRpcClient(target string, config *Config) *RpcClient { return &RpcClient{target, config} }

func (c *RpcClient) RunAgentToManagerRpc() error {
	conn, err := grpc.Dial(c.target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(c.target, err)
		return fmt.Errorf("failed to dial %s: %w", c.target, err)
	}
	defer conn.Close()

	managerClient := manager.NewGuardLinkManagerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := managerClient.NodeRegistration(ctx, &manager.NodeRegistrationRequest{
		Ip:              c.Config.conf.RegisterIP,
		InternetAddress: c.Config.conf.InternetAddress,
		NodeId:          c.Config.conf.AgentID,

		Protocol: func() manager.Protocol {
			if c.Config.conf.Protocol == "wss" {
				return manager.Protocol_WSS
			}

			return manager.Protocol_GRPC
		}(),
		WssPath: c.Config.WSPath,
	})

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to register %s: %w", c.target, err)
	}

	c.Config.conf.AgentID = response.NodeId
	c.Config.conf.WriteConf()
	log.Printf("NodeRegistration %s \n", response.NodeId)

	c.Config.SetJWTAESKey(response.AesKey)

	// 保持心跳
	for {
		time.Sleep(time.Second * 10)
		response, err := managerClient.NodeRegistration(context.TODO(), &manager.NodeRegistrationRequest{
			Ip:              c.Config.conf.RegisterIP,
			InternetAddress: c.Config.conf.InternetAddress,
			NodeId:          c.Config.conf.AgentID,

			Protocol: func() manager.Protocol {
				if c.Config.conf.Protocol == "wss" {
					return manager.Protocol_WSS
				}

				return manager.Protocol_GRPC
			}(),
			WssPath: c.Config.WSPath,
		})
		if err != nil {
			log.Println(err)
			continue
		}

		c.Config.SetJWTAESKey(response.AesKey)
	}
}
