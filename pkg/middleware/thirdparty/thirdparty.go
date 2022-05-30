package thirdparty

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/thirdlogingateway"
	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"
	crud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/thirdparty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Create(ctx context.Context, in *npool.ThirdParty) (*npool.ThirdParty, error) {
	schema, err := crud.New(context.Background(), nil)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	exist, err := schema.ExistConds(context.Background(),
		cruder.NewConds().WithCond(constant.ThirdPartyFieldBrandName, cruder.EQ, in.GetBrandName()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if exist {
		return nil, fmt.Errorf("brand name already exists")
	}
	schema, err = crud.New(context.Background(), nil)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp, err := schema.Create(context.Background(), in)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
