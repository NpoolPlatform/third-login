//nolint:dupl
package api

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/thirdlogingateway"
	"github.com/NpoolPlatform/third-login-gateway/pkg/auth"
	crud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/thirdparty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func checkThirdPartyInfo(info *npool.ThirdParty) error {
	if info.GetBrandName() == "" {
		logger.Sugar().Error("brand name is empty")
		return fmt.Errorf("app key empty")
	}

	if info.GetLogo() == "" {
		logger.Sugar().Error("logo is empty")
		return fmt.Errorf("logo empty")
	}

	if _, ok := auth.ThirdMap[info.GetDomain()]; ok {
		logger.Sugar().Error("unsupported login method")
		return fmt.Errorf("unsupported login method")
	}
	return nil
}

func (s *Server) CreateThirdParty(ctx context.Context, in *npool.CreateThirdPartyRequest) (*npool.CreateThirdPartyResponse, error) {
	err := checkThirdPartyInfo(in.GetInfo())
	if err != nil {
		return &npool.CreateThirdPartyResponse{}, status.Error(codes.Internal, err.Error())
	}
	schema, err := crud.New(context.Background(), nil)
	if err != nil {
		return &npool.CreateThirdPartyResponse{}, status.Error(codes.Internal, err.Error())
	}
	info, err := schema.Create(context.Background(), in.GetInfo())
	if err != nil {
		return &npool.CreateThirdPartyResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.CreateThirdPartyResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateThirdParty(ctx context.Context, in *npool.UpdateThirdPartyRequest) (*npool.UpdateThirdPartyResponse, error) {
	err := checkThirdPartyInfo(in.GetInfo())
	if err != nil {
		return &npool.UpdateThirdPartyResponse{}, status.Error(codes.Internal, err.Error())
	}
	schema, err := crud.New(context.Background(), nil)
	if err != nil {
		return &npool.UpdateThirdPartyResponse{}, status.Error(codes.Internal, err.Error())
	}
	info, err := schema.Update(context.Background(), in.GetInfo())
	if err != nil {
		return &npool.UpdateThirdPartyResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.UpdateThirdPartyResponse{
		Info: info,
	}, nil
}

func (s *Server) GetThirdParties(ctx context.Context, in *npool.GetThirdPartiesRequest) (*npool.GetThirdPartiesResponse, error) {
	schema, err := crud.New(context.Background(), nil)
	if err != nil {
		return &npool.GetThirdPartiesResponse{}, status.Error(codes.Internal, err.Error())
	}
	infos, _, err := schema.Rows(context.Background(), cruder.NewConds(), 0, 0)
	if err != nil {
		return &npool.GetThirdPartiesResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.GetThirdPartiesResponse{
		Infos: infos,
	}, nil
}
