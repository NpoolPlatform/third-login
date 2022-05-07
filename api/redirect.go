package api

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/third-login-gateway"
)

func (s *Server) ThirdRedirect(ctx context.Context, in *npool.ThirdLoginRequest) (*npool.ThirdLoginResponse, error) {
	return nil, nil
}
