package api

import (
	"context"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/third-login-gateway"
	mw "github.com/NpoolPlatform/third-login-gateway/pkg/middleware/platform"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetPlatformsByApp(ctx context.Context, in *npool.GetPlatformsByAppRequest) (*npool.GetPlatformsByAppResponse, error) {
	if _, err := uuid.Parse(in.AppID); err != nil {
		logger.Sugar().Errorf("invalid request app id: %v", err)
		return &npool.GetPlatformsByAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	resp, err := mw.GetPlatforms(ctx, in)
	if err != nil {
		logger.Sugar().Errorw("get platforms error: %v", err)
		return &npool.GetPlatformsByAppResponse{}, status.Error(codes.Internal, "internal server error")
	}
	return resp, nil
}
