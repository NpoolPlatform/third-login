package auth

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/thirdlogingateway"
	oauth "github.com/NpoolPlatform/third-login-gateway/pkg/auth"
	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"
	authcrud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/auth"
	thirdpartycrud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/thirdparty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetAuths(ctx context.Context, appID string) ([]*npool.Auth, error) {
	authSchema, err := authcrud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	authInfos, _, err := authSchema.Rows(ctx, cruder.NewConds().WithCond(constant.AuthFieldAppID, cruder.EQ, appID), 0, 0)
	if err != nil {
		logger.Sugar().Errorf("fail get auth: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	thirdPartySchema, err := thirdpartycrud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var authList []*npool.Auth
	for _, val := range authInfos {
		conf := &oauth.Config{ClientID: val.GetAppKey(), ClientSecret: val.GetAppSecret(), RedirectURL: val.GetRedirectURL()}
		thirdPartyInfo, err := thirdPartySchema.Row(ctx, uuid.MustParse(val.GetThirdPartyID()))
		if err != nil {
			logger.Sugar().Errorf("fail get auth: %v", err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		platform, ok := oauth.ThirdMap[thirdPartyInfo.GetDomain()]
		if !ok {
			return nil, fmt.Errorf("login method does not exist")
		}
		thirdMethod := oauth.NewContext(platform)
		url, err := thirdMethod.GetRedirectURL(conf)
		if err != nil {
			return nil, err
		}
		authList = append(authList, &npool.Auth{
			AppID:        val.GetAppID(),
			ThirdPartyID: val.GetThirdPartyID(),
			AuthURL:      url,
		})
	}
	return authList, nil
}
