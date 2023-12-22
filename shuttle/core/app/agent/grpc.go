package agent

import (
	"net"
)

type GrpcServer struct {
	net.Listener
}

// newGRPCWrapper wraps server, connection and rpc server
func newgRpcServer(conf *Config) (*GrpcServer, error) {
	listener, err := net.Listen("tcp", conf.Listen)
	if err != nil {
		return nil, err
	}
	return &GrpcServer{listener}, nil
}

func (s *GrpcServer) Read(b []byte) (n int, err error) {
	return 0, err
}

func (s *GrpcServer) Write(b []byte) (n int, err error) {
	return 0, err
}

func (s *GrpcServer) Serve() error {
	//srv := rpc.NewServer()
	//// todo: 注册rpc服务
	//// srv.RegisterService()
	//return srv.Serve(s.Listener)
	return nil
}

// func streamInterceptor(username string, password string) func(srv interface{}, stream rpc.ServerStream, info *rpc.StreamServerInfo, handler rpc.StreamHandler) error {
// 	return func(srv interface{}, stream rpc.ServerStream, info *rpc.StreamServerInfo, handler rpc.StreamHandler) error {
// 		if err := authorize(username, password)(stream.Context()); err != nil {
// 			return err
// 		}
// 		return handler(srv, stream)
// 	}
// }

// func unaryInterceptor(username string, password string) func(ctx context.Context, req interface{}, info *rpc.UnaryServerInfo, handler rpc.UnaryHandler) (interface{}, error) {
// 	return func(ctx context.Context, req interface{}, info *rpc.UnaryServerInfo, handler rpc.UnaryHandler) (interface{}, error) {
// 		if err := authorize(username, password)(ctx); err != nil {
// 			return nil, err
// 		}
// 		return handler(ctx, req)
// 	}
// }

// func authorize(username string, password string) func(ctx context.Context) error {
// 	return func(ctx context.Context) error {
// 		if md, ok := metadata.FromIncomingContext(ctx); ok {
// 			if len(md["username"]) > 0 && md["username"][0] == username &&
// 				len(md["password"]) > 0 && md["password"][0] == password {
// 				return nil
// 			}
// 			return fmt.Errorf("AccessDeniedErr")
// 		}
// 		return fmt.Errorf("EmptyMetadataErr")
// 	}
// }
