package rpc

import (
	"context"

	"github.com/rs/xid"
	"google.dev/google/shuttle/core/app/manager/conf"
	"google.dev/google/shuttle/core/app/manager/pkg/models"
	"google.dev/google/shuttle/core/app/manager/storage"
	"google.dev/google/shuttle/proto/manager"
)

type GRPCServer struct {
	manager.UnimplementedGuardLinkManagerServer
	Storage storage.Interface
}

func NewGRPCServer(storage storage.Interface) *GRPCServer {
	return &GRPCServer{Storage: storage}
}

func (g *GRPCServer) NodeRegistration(ctx context.Context, request *manager.NodeRegistrationRequest) (*manager.NodeRegistrationResponse, error) {
	if request.NodeId == "" {
		request.NodeId = xid.New().String()
	}

	var node = models.Node{
		NodeID:          request.NodeId,
		IP:              request.Ip,
		InternetAddress: request.InternetAddress,
		Protocol:        request.Protocol,
		WssPath:         request.WssPath,
	}

	err := g.Storage.DB().Model(&models.Node{}).Where("node_id = ?", request.NodeId).Attrs(&models.Node{
		IP:              request.Ip,
		InternetAddress: request.InternetAddress,
		Protocol:        request.Protocol,
		WssPath:         request.WssPath,
	}).FirstOrCreate(&node).Error
	if err != nil {
		return nil, err
	}

	return &manager.NodeRegistrationResponse{
		NodeId: request.NodeId,
		AesKey: conf.CONFIG.AgentAesKey,
	}, nil
}

func (g *GRPCServer) TrafficReport(ctx context.Context, request *manager.TrafficReportRequest) (*manager.TrafficReportResponse, error) {
	//TODO implement me
	panic("implement me")
}
