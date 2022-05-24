//nolint:dupl
package api

import (
	"context"
	thirdlgcrud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/thirdauth"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/thirdlogingateway"
	mw "github.com/NpoolPlatform/third-login-gateway/pkg/middleware/thirdauth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAuths(ctx context.Context, in *npool.GetAuthsRequest) (*npool.GetAuthsResponse, error) {
	if _, err := uuid.Parse(in.AppID); err != nil {
		logger.Sugar().Errorf("invalid request app id: %v", err)
		return &npool.GetAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}
	resp, err := mw.GetAuths(ctx, in)
	if err != nil {
		logger.Sugar().Errorw("get third auth error: %v", err)
		return &npool.GetAuthsResponse{}, status.Error(codes.Internal, "internal server error")
	}
	return resp, nil
}

func (s *Server) GetAuthsByApp(ctx context.Context, in *npool.GetAuthsByAppRequest) (*npool.GetAuthsByAppResponse, error) {
	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorf("invalid request target app id: %v", err)
		return &npool.GetAuthsByAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	resp, err := mw.GetAuthsByApp(ctx, in)
	if err != nil {
		logger.Sugar().Errorw("get third auth error: %v", err)
		return &npool.GetAuthsByAppResponse{}, status.Error(codes.Internal, "internal server error")
	}
	return resp, nil
}

func (s *Server) CreateAuth(ctx context.Context, in *npool.CreateAuthRequest) (*npool.CreateAuthResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorf("invalid request app id: %v", err)
		return &npool.CreateAuthResponse{}, status.Error(codes.Internal, err.Error())
	}
	schema, err := thirdlgcrud.New(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	info, err := schema.Create(context.Background(), in.GetInfo())
	if err != nil {
		return nil, err
	}
	return &npool.CreateAuthResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAuths(ctx context.Context, in *npool.CreateAuthsRequest) (*npool.CreateAuthsResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorf("invalid request app id: %v", err)
		return &npool.CreateAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}
	schema, err := thirdlgcrud.New(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	info, err := schema.CreateBulk(context.Background(), in.GetInfos())
	if err != nil {
		return nil, err
	}
	return &npool.CreateAuthsResponse{
		Infos: info,
	}, nil
}

func (s *Server) CreateAppAuth(ctx context.Context, in *npool.CreateAppAuthRequest) (*npool.CreateAppAuthResponse, error) {
	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorf("invalid request target app id: %v", err)
		return &npool.CreateAppAuthResponse{}, status.Error(codes.Internal, err.Error())
	}
	schema, err := thirdlgcrud.New(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	info, err := schema.Create(context.Background(), in.GetInfo())
	if err != nil {
		return nil, err
	}
	return &npool.CreateAppAuthResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAppAuths(ctx context.Context, in *npool.CreateAppAuthsRequest) (*npool.CreateAppAuthsResponse, error) {
	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorf("invalid request target app id: %v", err)
		return &npool.CreateAppAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}
	schema, err := thirdlgcrud.New(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	info, err := schema.CreateBulk(context.Background(), in.GetInfos())
	if err != nil {
		return nil, err
	}
	return &npool.CreateAppAuthsResponse{
		Infos: info,
	}, nil
}
