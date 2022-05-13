package thirdauth

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/third-login-gateway"
	oauth "github.com/NpoolPlatform/third-login-gateway/pkg/auth"
	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"
	crud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/thirdauth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetThirdAuth(ctx context.Context, in *npool.GetThirdAuthByAppRequest) (*npool.GetThirdAuthByAppResponse, error) {
	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.GetThirdAuthByAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	infos, _, err := schema.Rows(ctx, cruder.NewConds().WithCond(constant.ThirdAuthFieldAppID, cruder.EQ, in.GetAppID()), 0, 0)
	if err != nil {
		logger.Sugar().Errorf("fail get third auth: %v", err)
		return &npool.GetThirdAuthByAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	var authList []*npool.Auth

	for _, val := range infos {
		conf := &oauth.Config{ClientID: val.ThirdAppKey, ClientSecret: val.ThirdAppSecret, RedirectURL: val.RedirectUrl}
		platform, ok := oauth.ThirdMap[val.GetThird()]
		if !ok {
			return &npool.GetThirdAuthByAppResponse{}, fmt.Errorf("login method does not exist")
		}
		thirdMethod := oauth.NewContext(platform)
		url, err := thirdMethod.GetRedirectURL(conf)
		if err != nil {
			return &npool.GetThirdAuthByAppResponse{}, err
		}
		authList = append(authList, &npool.Auth{
			AuthUrl: url,
			LogoUrl: val.LogoUrl,
			Third:   val.Third,
		})
	}
	return &npool.GetThirdAuthByAppResponse{
		Infos: authList,
	}, nil
}
