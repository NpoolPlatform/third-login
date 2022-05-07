package platform

import (
	"context"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/third-login-gateway"
	oauth "github.com/NpoolPlatform/third-login-gateway/pkg/auth"
	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"
	crud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/platform"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetPlatforms(ctx context.Context, in *npool.GetPlatformsByAppRequest) (*npool.GetPlatformsByAppResponse, error) {

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.GetPlatformsByAppResponse{}, status.Error(codes.Internal, err.Error())
	}

	infos, _, err := schema.Rows(ctx, cruder.NewConds().WithCond(constant.PlatformFieldAppID, cruder.EQ, in.GetAppID()), 0, 0)
	if err != nil {
		logger.Sugar().Errorf("fail platform stocks: %v", err)
		return &npool.GetPlatformsByAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	var autlList []*npool.Auth
	for _, val := range infos {
		conf := &oauth.AuthConfig{ClientId: val.PlatformAppKey, ClientSecret: val.PlatformAppSecret, RedirectUrl: val.RedirectUrl}
		switch val.GetPlatform() {
		case constant.PlatformGitHub:
			githubAuth := oauth.NewAuthGitHub(conf)
			authUrl, err := githubAuth.GetRedirectUrl()
			if err != nil {
				return nil, err
			}
			autlList = append(autlList, &npool.Auth{
				AuthUrl: authUrl,
				LogoUrl: val.LogoUrl,
			})
			break
		case constant.PlatformGooGle:
			googleAuth := oauth.NewAuthGoogle(conf)
			authUrl, err := googleAuth.GetRedirectUrl()
			if err != nil {
				return nil, err
			}
			autlList = append(autlList, &npool.Auth{
				AuthUrl: authUrl,
				LogoUrl: val.LogoUrl,
			})
			break
		}
	}
	return &npool.GetPlatformsByAppResponse{
		Infos: autlList,
	}, nil
}
