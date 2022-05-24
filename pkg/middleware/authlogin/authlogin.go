package authlogin

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	npool "github.com/NpoolPlatform/message/npool/thirdlogingateway"
	oauth "github.com/NpoolPlatform/third-login-gateway/pkg/auth"
	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"
	crud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/thirdauth"
	grpc2 "github.com/NpoolPlatform/third-login-gateway/pkg/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Login(ctx context.Context, in *npool.LoginRequest) (*npool.LoginResponse, error) {
	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.LoginResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := schema.RowOnly(ctx, cruder.NewConds().
		WithCond(constant.ThirdAuthFieldAppID, cruder.EQ, in.GetAppID()).
		WithCond(constant.ThirdAuthFieldThird, cruder.EQ, in.GetThird()))
	if err != nil {
		logger.Sugar().Errorf("fail get third auth: %v", err)
		return &npool.LoginResponse{}, status.Error(codes.Internal, err.Error())
	}

	conf := &oauth.Config{ClientID: info.ThirdAppKey, ClientSecret: info.ThirdAppSecret, RedirectURL: info.RedirectUrl}
	third, ok := oauth.ThirdMap[in.GetThird()]
	if !ok {
		return &npool.LoginResponse{}, fmt.Errorf("login method does not exist")
	}
	thirdMethod := oauth.NewContext(third)
	thirdUser, err := thirdMethod.GetUserInfo(ctx, in.GetCode(), conf)
	if err != nil {
		return &npool.LoginResponse{}, err
	}
	thirdUser.AppID = in.GetAppID()

	var tUser *appusermgrpb.AppUserThird
	tUser, err = grpc2.GetAppUserThirdByAppThird(ctx, &appusermgrpb.GetAppUserThirdByAppThirdRequest{
		AppID:       thirdUser.AppID,
		ThirdID:     thirdUser.ThirdID,
		ThirdUserID: thirdUser.ThirdUserID,
	})
	if err != nil {
		return &npool.LoginResponse{}, err
	}
	if tUser == nil {
		user, err := grpc2.CreateAppUserWithThird(ctx, &appusermgrpb.CreateAppUserWithThirdRequest{
			User: &appusermgrpb.AppUser{
				AppID: thirdUser.AppID,
			},
			Third: thirdUser,
		})
		if err != nil {
			return &npool.LoginResponse{}, err
		}
		if user == nil {
			return &npool.LoginResponse{}, fmt.Errorf("fail createa app user with third")
		}
		thirdUser.ID = user.GetID()
	}

	userInfo, err := grpc2.GetAppUserInfo(ctx, &appusermgrpb.GetAppUserInfoRequest{
		ID: thirdUser.ID,
	})
	if err != nil {
		return &npool.LoginResponse{}, err
	}
	return &npool.LoginResponse{
		Info: userInfo,
	}, err
}
