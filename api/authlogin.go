package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	mw "github.com/NpoolPlatform/third-login-gateway/pkg/middleware/authlogin"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	npool "github.com/NpoolPlatform/message/npool/third-login-gateway"
)

func (s *Server) AuthLogin(ctx context.Context, in *npool.AuthLoginRequest) (*npool.AuthLoginResponse, error) {
	if in.GetCode() == "" {
		logger.Sugar().Error("auth login error code is empty")
		return nil, status.Error(codes.InvalidArgument, "Code empty")
	}

	if in.GetThird() == "" {
		logger.Sugar().Error("auth login error third is empty")
		return nil, status.Error(codes.InvalidArgument, "Platform empty")
	}

	if _, err := uuid.Parse(in.AppID); err != nil {
		logger.Sugar().Errorf("invalid request app id: %v", err)
		return &npool.AuthLoginResponse{}, status.Error(codes.Internal, err.Error())
	}
	resp, err := mw.AuthLogin(ctx, in)
	if err != nil {
		logger.Sugar().Errorw("auth login error: %v", err)
		return &npool.AuthLoginResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
