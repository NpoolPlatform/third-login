package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	mw "github.com/NpoolPlatform/third-login-gateway/pkg/middleware/authlogin"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	npool "github.com/NpoolPlatform/message/npool/thirdlogingateway"
)

func (s *Server) Login(ctx context.Context, in *npool.LoginRequest) (*npool.LoginResponse, error) {
	if in.GetCode() == "" {
		logger.Sugar().Error("auth login error code is empty")
		return nil, status.Error(codes.InvalidArgument, "code empty")
	}

	if in.GetThird() == "" {
		logger.Sugar().Error("auth login error third is empty")
		return nil, status.Error(codes.InvalidArgument, "third empty")
	}

	if _, err := uuid.Parse(in.AppID); err != nil {
		logger.Sugar().Errorf("invalid request app id: %v", err)
		return &npool.LoginResponse{}, status.Error(codes.Internal, err.Error())
	}
	resp, err := mw.Login(ctx, in)
	if err != nil {
		logger.Sugar().Errorw("auth login error: %v", err)
		return &npool.LoginResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
