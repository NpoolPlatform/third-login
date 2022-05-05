package api

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/third-login-gateway"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	npool.UnimplementedThirdLoginGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	npool.RegisterThirdLoginGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return npool.RegisterThirdLoginGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
