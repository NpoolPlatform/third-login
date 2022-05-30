//nolint:dupl
package api

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/thirdlogingateway"
	crud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/auth"
	mw "github.com/NpoolPlatform/third-login-gateway/pkg/middleware/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAuths(ctx context.Context, in *npool.GetAuthsRequest) (*npool.GetAuthsResponse, error) {
	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorf("invalid request app id: %v", err)
		return &npool.GetAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}
	resp, err := mw.GetAuths(ctx, in.GetAppID())
	if err != nil {
		logger.Sugar().Errorw("get third auth error: %v", err)
		return &npool.GetAuthsResponse{}, status.Error(codes.Internal, "internal server error")
	}
	return &npool.GetAuthsResponse{
		Infos: resp,
	}, nil
}

func (s *Server) GetAppAuths(ctx context.Context, in *npool.GetAppAuthsRequest) (*npool.GetAppAuthsResponse, error) {
	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorf("invalid request target app id: %v", err)
		return &npool.GetAppAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}
	resp, err := mw.GetAuths(ctx, in.GetTargetAppID())
	if err != nil {
		logger.Sugar().Errorw("get third auth error: %v", err)
		return &npool.GetAppAuthsResponse{}, status.Error(codes.Internal, "internal server error")
	}
	return &npool.GetAppAuthsResponse{
		Infos: resp,
	}, nil
}

func checkAuthInfo(info *npool.Auth) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		logger.Sugar().Errorf("invalid request app id: %v", err)
		return err
	}

	if _, err := uuid.Parse(info.GetThirdPartyID()); err != nil {
		logger.Sugar().Errorf("invalid request third party id: %v", err)
		return err
	}

	if info.GetAppKey() == "" {
		logger.Sugar().Error("app key is empty")
		return fmt.Errorf("app key empty")
	}

	if info.GetAppSecret() == "" {
		logger.Sugar().Error("app secret is empty")
		return fmt.Errorf("app key empty")
	}

	if info.GetRedirectURL() == "" {
		logger.Sugar().Error("redirect url is empty")
		return fmt.Errorf("app key empty")
	}
	return nil
}

func (s *Server) CreateAuth(ctx context.Context, in *npool.CreateAuthRequest) (*npool.CreateAuthResponse, error) {
	err := checkAuthInfo(in.GetInfo())
	if err != nil {
		return &npool.CreateAuthResponse{}, status.Error(codes.Internal, err.Error())
	}
	schema, err := crud.New(context.Background(), nil)
	if err != nil {
		return &npool.CreateAuthResponse{}, status.Error(codes.Internal, err.Error())
	}
	info, err := schema.Create(context.Background(), in.GetInfo())
	if err != nil {
		return &npool.CreateAuthResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.CreateAuthResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAuths(ctx context.Context, in *npool.CreateAuthsRequest) (*npool.CreateAuthsResponse, error) {
	for _, val := range in.GetInfos() {
		err := checkAuthInfo(val)
		if err != nil {
			return &npool.CreateAuthsResponse{}, status.Error(codes.Internal, err.Error())
		}
	}
	schema, err := crud.New(context.Background(), nil)
	if err != nil {
		return &npool.CreateAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}
	info, err := schema.CreateBulk(context.Background(), in.GetInfos())
	if err != nil {
		return &npool.CreateAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.CreateAuthsResponse{
		Infos: info,
	}, nil
}

func (s *Server) CreateAppAuth(ctx context.Context, in *npool.CreateAppAuthRequest) (*npool.CreateAppAuthResponse, error) {
	err := checkAuthInfo(in.GetInfo())
	if err != nil {
		return &npool.CreateAppAuthResponse{}, status.Error(codes.Internal, err.Error())
	}
	schema, err := crud.New(context.Background(), nil)
	if err != nil {
		return &npool.CreateAppAuthResponse{}, status.Error(codes.Internal, err.Error())
	}
	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		logger.Sugar().Errorf("invalid request target app id: %v", err)
		return &npool.CreateAppAuthResponse{}, status.Error(codes.Internal, err.Error())
	}
	in.GetInfo().AppID = in.GetTargetAppID()
	info, err := schema.Create(context.Background(), in.GetInfo())
	if err != nil {
		return &npool.CreateAppAuthResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.CreateAppAuthResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateAppAuths(ctx context.Context, in *npool.CreateAppAuthsRequest) (*npool.CreateAppAuthsResponse, error) {
	for key, val := range in.GetInfos() {
		err := checkAuthInfo(val)
		if err != nil {
			return &npool.CreateAppAuthsResponse{}, status.Error(codes.Internal, err.Error())
		}
		in.GetInfos()[key].AppID = in.GetTargetAppID()
	}
	schema, err := crud.New(context.Background(), nil)
	if err != nil {
		return &npool.CreateAppAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}
	info, err := schema.CreateBulk(context.Background(), in.GetInfos())
	if err != nil {
		return &npool.CreateAppAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.CreateAppAuthsResponse{
		Infos: info,
	}, nil
}
