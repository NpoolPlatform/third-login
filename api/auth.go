//nolint:dupl
package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/thirdlogingateway"
	crud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/auth"
	mw "github.com/NpoolPlatform/third-login-gateway/pkg/middleware/auth"
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

func (s *Server) GetAppAuths(ctx context.Context, in *npool.GetAppAuthsRequest) (*npool.GetAppAuthsResponse, error) {
	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorf("invalid request target app id: %v", err)
		return &npool.GetAppAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}
	resp, err := mw.GetAppAuths(ctx, in)
	if err != nil {
		logger.Sugar().Errorw("get third auth error: %v", err)
		return &npool.GetAppAuthsResponse{}, status.Error(codes.Internal, "internal server error")
	}
	return resp, nil
}

func checkAuthInfo(info *npool.Auth) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		logger.Sugar().Errorf("invalid request app id: %v", err)
		return status.Error(codes.Internal, err.Error())
	}

	if _, err := uuid.Parse(info.GetThirdPartyID()); err != nil {
		logger.Sugar().Errorf("invalid request third party id: %v", err)
		return status.Error(codes.Internal, err.Error())
	}

	if info.GetAppKey() == "" {
		logger.Sugar().Error("app key is empty")
		return status.Error(codes.Internal, "app key empty")
	}

	if info.GetAppSecret() == "" {
		logger.Sugar().Error("app secret is empty")
		return status.Error(codes.Internal, "app secret empty")
	}

	if info.GetRedirectURL() == "" {
		logger.Sugar().Error("redirect url is empty")
		return status.Error(codes.Internal, "redirect url empty")
	}
	return nil
}

func (s *Server) CreateAuth(ctx context.Context, in *npool.CreateAuthRequest) (*npool.CreateAuthResponse, error) {
	err := checkAuthInfo(in.GetInfo())
	if err != nil {
		return &npool.CreateAuthResponse{}, err
	}
	schema, err := crud.New(context.Background(), nil)
	if err != nil {
		return &npool.CreateAuthResponse{}, err
	}
	info, err := schema.Create(context.Background(), in.GetInfo())
	if err != nil {
		return &npool.CreateAuthResponse{}, err
	}
	return &npool.CreateAuthResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAuths(ctx context.Context, in *npool.CreateAuthsRequest) (*npool.CreateAuthsResponse, error) {
	for _, val := range in.GetInfos() {
		err := checkAuthInfo(val)
		if err != nil {
			return &npool.CreateAuthsResponse{}, err
		}
	}
	schema, err := crud.New(context.Background(), nil)
	if err != nil {
		return &npool.CreateAuthsResponse{}, err
	}
	info, err := schema.CreateBulk(context.Background(), in.GetInfos())
	if err != nil {
		return &npool.CreateAuthsResponse{}, err
	}
	return &npool.CreateAuthsResponse{
		Infos: info,
	}, nil
}

func (s *Server) CreateAppAuth(ctx context.Context, in *npool.CreateAppAuthRequest) (*npool.CreateAppAuthResponse, error) {
	err := checkAuthInfo(in.GetInfo())
	if err != nil {
		return &npool.CreateAppAuthResponse{}, err
	}
	schema, err := crud.New(context.Background(), nil)
	if err != nil {
		return &npool.CreateAppAuthResponse{}, err
	}
	info, err := schema.Create(context.Background(), in.GetInfo())
	if err != nil {
		return &npool.CreateAppAuthResponse{}, err
	}
	return &npool.CreateAppAuthResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAppAuths(ctx context.Context, in *npool.CreateAppAuthsRequest) (*npool.CreateAppAuthsResponse, error) {
	for _, val := range in.GetInfos() {
		err := checkAuthInfo(val)
		if err != nil {
			return &npool.CreateAppAuthsResponse{}, err
		}
	}
	schema, err := crud.New(context.Background(), nil)
	if err != nil {
		return &npool.CreateAppAuthsResponse{}, err
	}
	info, err := schema.CreateBulk(context.Background(), in.GetInfos())
	if err != nil {
		return &npool.CreateAppAuthsResponse{}, err
	}
	return &npool.CreateAppAuthsResponse{
		Infos: info,
	}, nil
}
