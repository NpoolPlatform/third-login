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

func GetAuths(ctx context.Context, in *npool.GetAuthsRequest) (*npool.GetAuthsResponse, error) {
	authSchema, err := authcrud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.GetAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}

	authInfos, _, err := authSchema.Rows(ctx, cruder.NewConds().WithCond(constant.AuthFieldAppID, cruder.EQ, in.GetAppID()), 0, 0)
	if err != nil {
		logger.Sugar().Errorf("fail get auth: %v", err)
		return &npool.GetAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}

	thirdPartySchema, err := thirdpartycrud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.GetAuthsResponse{}, status.Error(codes.Internal, err.Error())
	}

	var authList []*npool.Auth
	for _, val := range authInfos {
		conf := &oauth.Config{ClientID: val.GetAppKey(), ClientSecret: val.GetAppSecret(), RedirectURL: val.GetRedirectURL()}
		thirdPartyInfo, err := thirdPartySchema.Row(ctx, uuid.MustParse(val.GetThirdPartyID()))
		if err != nil {
			logger.Sugar().Errorf("fail get auth: %v", err)
			return &npool.GetAuthsResponse{}, status.Error(codes.Internal, err.Error())
		}
		platform, ok := oauth.ThirdMap[thirdPartyInfo.GetDomain()]
		if !ok {
			return &npool.GetAuthsResponse{}, fmt.Errorf("login method does not exist")
		}
		thirdMethod := oauth.NewContext(platform)
		url, err := thirdMethod.GetRedirectURL(conf)
		if err != nil {
			return &npool.GetAuthsResponse{}, err
		}
		authList = append(authList, &npool.Auth{
			AppID:        val.GetAppID(),
			ThirdPartyID: val.GetThirdPartyID(),
			AuthURL:      url,
		})
	}
	return &npool.GetAuthsResponse{
		Infos: authList,
	}, nil
}

func GetAppAuths(ctx context.Context, in *npool.GetAppAuthsRequest) (*npool.GetAppAuthsResponse, error) {
	resp, err := GetAuths(ctx, &npool.GetAuthsRequest{AppID: in.GetTargetAppID()})
	if err != nil {
		return nil, err
	}
	return &npool.GetAppAuthsResponse{
		Infos: resp.GetInfos(),
	}, err
}
