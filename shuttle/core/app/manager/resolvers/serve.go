package resolvers

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.dev/google/shuttle/core/app/manager/conf"
	"google.dev/google/shuttle/core/app/manager/generated"
	"google.dev/google/shuttle/core/app/manager/pkg/errs"
	"google.dev/google/shuttle/core/app/manager/pkg/models"
	"google.dev/google/shuttle/core/app/manager/utils"
	"google.dev/google/shuttle/pkg"
	"google.dev/google/shuttle/proto/manager"
	"google.dev/google/socks5_discovery/proto"
)

func (r *queryResolver) Combos(ctx context.Context) (*generated.Combos, error) {
	authJWT, err := utils.GetUserInformationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var combos []models.Combo
	err = r.Storage.DB().Model(&models.Combo{}).
		Where("app_id = ?", authJWT.AppID).Order("sort asc").Find(&combos).Error
	if err != nil {
		return nil, err
	}

	var result generated.Combos
	for _, v := range combos {
		result.Combos = append(result.Combos, generated.Combo{
			ComboID:  v.ComboID,
			Describe: v.Describe,
			Traffic:  v.Traffic,
			Day:      v.Day,
			Amount:   v.Amount,
		})
	}

	return &result, nil
}

func (r *queryResolver) NodeToken(ctx context.Context, nodeID string, mountToken *string) (*generated.NodeToken, error) {
	authJWT, err := utils.GetUserInformationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var appNodeMapping models.AppNodeMapping
	err = r.Storage.DB().Model(&models.AppNodeMapping{}).
		Where("app_id = ?", authJWT.AppID).
		Where("node_id = ?", nodeID).First(&appNodeMapping).Error
	if err != nil {
		return nil, errs.BadRequest
	}

	// TODO: 验证节点 free or vip
	var node models.Node
	err = r.Storage.DB().Model(&models.Node{}).Where("node_id = ?", nodeID).First(&node).Error
	if err != nil {
		return nil, err
	}

	var vt = &pkg.VerifyToken{
		UserJWT:    authJWT.Token,
		NodeID:     nodeID,
		Expiration: time.Now().Add(time.Hour * 5).Unix(),
	}

	if mountToken != nil {
		var pt = &proto.Socks5{}
		pt.FromToken(*mountToken, conf.CONFIG.AgentAesKey)
		vt.MountSocks5 = pkg.MountSocks5{
			Address:  pt.Address,
			Username: pt.Username,
			Password: pt.Password,
		}
	}

	return &generated.NodeToken{
		InternetAddress: node.InternetAddress,
		NodeProtocol: func() generated.NodeProtocol {
			if node.Protocol == manager.Protocol_WSS {
				return generated.NodeProtocolWss
			}

			return generated.NodeProtocolRPC
		}(),
		WsPath: node.WssPath,
		Token:  vt.ToToken(conf.CONFIG.AgentAesKey),
	}, nil
}

// Nodes ...
func (r *queryResolver) Nodes(ctx context.Context) (*generated.Nodes, error) {
	authJWT, err := utils.GetUserInformationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var anms []models.AppNodeMapping
	err = r.Storage.DB().Model(&models.AppNodeMapping{}).
		Where("app_id = ?", authJWT.AppID).
		Order("sort asc").Find(&anms).Error
	if err != nil {
		return nil, err
	}

	var nodeIds []string
	for _, v := range anms {
		nodeIds = append(nodeIds, v.NodeID)
	}

	var nodes []models.Node
	err = r.Storage.DB().Model(&models.Node{}).Where("node_id in ?", nodeIds).Find(&nodes).Error
	if err != nil {
		return nil, err
	}

	for i, v := range anms {
		for ii, vv := range nodes {
			if vv.NodeID == v.NodeID {
				anms[i].Node = &nodes[ii]
			}
		}
	}

	var mountNodeID string
	var result generated.Nodes
	for _, v := range anms {
		result.Nodes = append(result.Nodes, generated.NodeItem{
			NodeID:   v.NodeID,
			NodeName: v.Node.NodeName,
			Country:  v.Node.Country,
			Describe: v.Node.Describe,
			Free:     v.Free,
		})

		if v.Node.MountSupport {
			mountNodeID = v.NodeID
		}
	}

	discovery, err := r.socks5DiscoveryClient.Discovery(context.TODO(), &proto.DiscoveryRequest{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for i, v := range discovery.Socks5S {
		result.MountNodes = append(result.MountNodes, generated.MountNode{
			NodeID:     mountNodeID,
			NodeName:   fmt.Sprintf("%s-%d", v.Country, i),
			Country:    v.Country,
			Describe:   fmt.Sprintf("Cloud node, speed unstable, latency %d", v.Delay/10),
			MountToken: v.ToToken(conf.CONFIG.AgentAesKey),
		})
	}

	return &result, nil
}
