package authlogin

import (
	"context"
	"fmt"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	npool "github.com/NpoolPlatform/message/npool/third-login-gateway"
	oauth "github.com/NpoolPlatform/third-login-gateway/pkg/auth"
	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"
	crud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/thirdauth"
	grpc2 "github.com/NpoolPlatform/third-login-gateway/pkg/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AuthLogin(ctx context.Context, in *npool.AuthLoginRequest) (*npool.AuthLoginResponse, error) {
	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.AuthLoginResponse{}, status.Error(codes.Internal, err.Error())
	}

	infos, _, err := schema.Rows(ctx, cruder.NewConds().
		WithCond(constant.ThirdAuthFieldAppID, cruder.EQ, in.GetAppID()).
		WithCond(constant.ThirdAuthFieldThird, cruder.EQ, in.GetThird()), 0, 100)
	if err != nil {
		logger.Sugar().Errorf("fail get third auth: %v", err)
		return &npool.AuthLoginResponse{}, status.Error(codes.Internal, err.Error())
	}
	if len(infos) == 0 {
		return &npool.AuthLoginResponse{}, status.Error(codes.Internal, "not find third auth")
	}

	conf := &oauth.Config{ClientID: infos[0].ThirdAppKey, ClientSecret: infos[0].ThirdAppSecret, RedirectURL: infos[0].RedirectUrl}
	platform, ok := oauth.ThirdMap[in.GetThird()]
	if !ok {
		return &npool.AuthLoginResponse{}, fmt.Errorf("login method does not exist")
	}
	thirdMethod := oauth.NewContext(platform)
	thirdUser, err := thirdMethod.GetUserInfo(ctx, in.GetCode(), conf)
	if err != nil {
		return &npool.AuthLoginResponse{}, err
	}

	var tUser *appusermgrpb.AppUserThird
	tUser, err = grpc2.GetAppUserThirdByAppThird(ctx, &appusermgrpb.GetAppUserThirdByAppThirdRequest{
		AppID:       thirdUser.AppID,
		ThirdID:     thirdUser.ThirdId,
		ThirdUserID: thirdUser.ThirdUserId,
	})
	if err != nil {
		return &npool.AuthLoginResponse{}, err
	}

	if tUser == nil {
		user, err := grpc2.CreateAppUserWithThird(ctx, &appusermgrpb.CreateAppUserWithThirdRequest{
			User: &appusermgrpb.AppUser{
				AppID: thirdUser.AppID,
			},
			Third: thirdUser,
		})
		if err != nil {
			return &npool.AuthLoginResponse{}, err
		}
		tUser.UserID = user.ID
	}

	userInfo, err := grpc2.GetAppUserInfo(ctx, &appusermgrpb.GetAppUserInfoRequest{
		ID: tUser.UserID,
	})
	if err != nil {
		return &npool.AuthLoginResponse{}, err
	}
	return &npool.AuthLoginResponse{
		Info: userInfo,
	}, err
}
