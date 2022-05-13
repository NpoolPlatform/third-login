package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/third-login-gateway"
	mw "github.com/NpoolPlatform/third-login-gateway/pkg/middleware/thirdauth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetThirdAuthByApp(ctx context.Context, in *npool.GetThirdAuthByAppRequest) (*npool.GetThirdAuthByAppResponse, error) {
	if _, err := uuid.Parse(in.AppID); err != nil {
		logger.Sugar().Errorf("invalid request app id: %v", err)
		return &npool.GetThirdAuthByAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	resp, err := mw.GetThirdAuth(ctx, in)
	if err != nil {
		logger.Sugar().Errorw("get third auth error: %v", err)
		return &npool.GetThirdAuthByAppResponse{}, status.Error(codes.Internal, "internal server error")
	}
	return resp, nil
}
