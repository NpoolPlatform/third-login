package platform

import (
	"context"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/third-login-gateway"
	oauth "github.com/NpoolPlatform/third-login-gateway/pkg/auth"
	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"
	crud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/platform"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetPlatformsAuth(ctx context.Context, in *npool.GetPlatformsByAppRequest) (*npool.GetPlatformsByAppResponse, error) {
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
		conf := &oauth.Config{ClientID: val.PlatformAppKey, ClientSecret: val.PlatformAppSecret, RedirectURL: val.RedirectUrl}
		switch val.GetPlatform() {
		case appusermgrconst.ThirdGithub:
			githubAuth := oauth.NewGitHubAuth(conf)
			authURL, err := githubAuth.GetRedirectURL()
			if err != nil {
				return nil, err
			}
			autlList = append(autlList, &npool.Auth{
				AuthUrl:  authURL,
				LogoUrl:  val.LogoUrl,
				Platform: val.Platform,
			})
			break //nolint
		case appusermgrconst.ThirdGoogle:
			googleAuth := oauth.NewGoogleAuth(conf)
			authURL, err := googleAuth.GetRedirectURL()
			if err != nil {
				return nil, err
			}
			autlList = append(autlList, &npool.Auth{
				AuthUrl:  authURL,
				LogoUrl:  val.LogoUrl,
				Platform: val.Platform,
			})
			break //nolint
		}
	}
	return &npool.GetPlatformsByAppResponse{
		Infos: autlList,
	}, nil
}
