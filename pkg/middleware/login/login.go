package login

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	npool "github.com/NpoolPlatform/message/npool/thirdlogingateway"
	oauth "github.com/NpoolPlatform/third-login-gateway/pkg/auth"
	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"
	authcrud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/auth"
	thirdpartycrud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/thirdparty"
	grpc2 "github.com/NpoolPlatform/third-login-gateway/pkg/grpc"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Login(ctx context.Context, in *npool.LoginRequest) (*npool.LoginResponse, error) {
	authSchema, err := authcrud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.LoginResponse{}, status.Error(codes.Internal, err.Error())
	}

	authInfo, err := authSchema.RowOnly(ctx, cruder.NewConds().
		WithCond(constant.AuthFieldAppID, cruder.EQ, in.GetAppID()).
		WithCond(constant.AuthFieldThirdPartyID, cruder.EQ, in.GetThirdPartyID()))
	if err != nil {
		logger.Sugar().Errorf("fail get third auth: %v", err)
		return &npool.LoginResponse{}, status.Error(codes.Internal, err.Error())
	}

	thirdPartySchema, err := thirdpartycrud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.LoginResponse{}, status.Error(codes.Internal, err.Error())
	}

	thirdPartyInfo, err := thirdPartySchema.Row(ctx, uuid.MustParse(authInfo.GetThirdPartyID()))
	if err != nil {
		logger.Sugar().Errorf("fail get third auth: %v", err)
		return &npool.LoginResponse{}, status.Error(codes.Internal, err.Error())
	}

	conf := &oauth.Config{ClientID: authInfo.AppKey, ClientSecret: authInfo.AppSecret, RedirectURL: authInfo.RedirectURL}
	third, ok := oauth.ThirdMap[thirdPartyInfo.GetDomain()]
	if !ok {
		return &npool.LoginResponse{}, fmt.Errorf("login method does not exist")
	}
	thirdMethod := oauth.NewContext(third)
	thirdUser, err := thirdMethod.GetUserInfo(ctx, in.GetCode(), conf)
	if err != nil {
		return &npool.LoginResponse{}, err
	}
	thirdUser.AppID = in.GetAppID()

	var tUser *appusermgrpb.AppUserThirdParty
	tUser, err = grpc2.GetAppUserThirdByAppThird(ctx, &appusermgrpb.GetAppUserThirdPartyByAppThirdPartyIDRequest{
		AppID:            thirdUser.AppID,
		ThirdPartyID:     thirdUser.ThirdPartyID,
		ThirdPartyUserID: thirdUser.ThirdPartyUserID,
	})
	if err != nil {
		return &npool.LoginResponse{}, err
	}
	if tUser == nil {
		user, err := grpc2.CreateAppUserWithThird(ctx, &appusermgrpb.CreateAppUserWithThirdPartyRequest{
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
