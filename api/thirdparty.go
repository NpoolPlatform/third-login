package api

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/thirdlogingateway"
	"github.com/NpoolPlatform/third-login-gateway/pkg/auth"
	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"
	crud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/thirdparty"
	mw "github.com/NpoolPlatform/third-login-gateway/pkg/middleware/thirdparty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func checkThirdPartyInfo(info *npool.ThirdParty) error {
	if info.GetBrandName() == "" {
		logger.Sugar().Error("brand name is empty")
		return fmt.Errorf("brand name is empty")
	}

	if info.GetLogo() == "" {
		logger.Sugar().Error("logo is empty")
		return fmt.Errorf("logo empty")
	}

	if _, ok := auth.ThirdMap[info.GetBrandName()]; !ok {
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
	resp, err := mw.Create(ctx, in.GetInfo())
	if err != nil {
		return &npool.CreateThirdPartyResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.CreateThirdPartyResponse{
		Info: resp,
	}, nil
}

func (s *Server) UpdateThirdParty(ctx context.Context, in *npool.UpdateThirdPartyRequest) (*npool.UpdateThirdPartyResponse, error) {
	err := checkThirdPartyInfo(in.GetInfo())
	if _, err := uuid.Parse(in.GetInfo().GetID()); err != nil {
		logger.Sugar().Errorf("invalid request app id: %v", err)
		return &npool.UpdateThirdPartyResponse{}, status.Error(codes.Internal, err.Error())
	}
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

func thirdPartyCondsToConds(conds cruder.FilterConds) (cruder.Conds, error) {
	newConds := cruder.NewConds()

	for k, v := range conds {
		switch v.Op {
		case cruder.EQ:
		case cruder.GT:
		case cruder.LT:
		case cruder.LIKE:
		default:
			return nil, fmt.Errorf("invalid filter condition op")
		}

		switch k {
		case constant.FieldID:
			fallthrough //nolint
		case constant.ThirdPartyFieldBrandName:
			newConds = newConds.WithCond(k, v.Op, v.Val.GetStringValue())
		default:
			return nil, fmt.Errorf("invalid third party field")
		}
	}

	return newConds, nil
}

func (s *Server) GetThirdPartyOnly(ctx context.Context, in *npool.GetThirdPartyOnlyRequest) (*npool.GetThirdPartyOnlyResponse, error) {
	conds, err := thirdPartyCondsToConds(in.GetConds())
	if err != nil {
		logger.Sugar().Errorf("invalid stock fields: %v", err)
		return &npool.GetThirdPartyOnlyResponse{}, status.Error(codes.Internal, err.Error())
	}
	schema, err := crud.New(context.Background(), nil)
	if err != nil {
		return &npool.GetThirdPartyOnlyResponse{}, status.Error(codes.Internal, err.Error())
	}
	info, err := schema.RowOnly(context.Background(), conds)
	if err != nil {
		return &npool.GetThirdPartyOnlyResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.GetThirdPartyOnlyResponse{
		Info: info,
	}, nil
}
