package authlogin

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/third-login-gateway"
	oauth "github.com/NpoolPlatform/third-login-gateway/pkg/auth"
	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"
	crud "github.com/NpoolPlatform/third-login-gateway/pkg/crud/platform"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	grpc2 "github.com/NpoolPlatform/third-login-gateway/pkg/grpc"
)

func AuthLogin(ctx context.Context, in *npool.AuthLoginRequest) (*npool.AuthLoginResponse, error) {
	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.AuthLoginResponse{}, status.Error(codes.Internal, err.Error())
	}

	infos, _, err := schema.Rows(ctx, cruder.NewConds().
		WithCond(constant.PlatformFieldAppID, cruder.EQ, in.GetAppID()).
		WithCond(constant.PlatformFieldPlatform, cruder.EQ, in.GetPlatform()), 0, 100)
	if err != nil {
		logger.Sugar().Errorf("fail get platform: %v", err)
		return &npool.AuthLoginResponse{}, status.Error(codes.Internal, err.Error())
	}
	if len(infos) == 0 {
		return &npool.AuthLoginResponse{}, status.Error(codes.Internal, "not find platform")
	}

	conf := &oauth.Config{ClientID: infos[0].PlatformAppKey, ClientSecret: infos[0].PlatformAppSecret, RedirectURL: infos[0].RedirectUrl}
	switch in.GetPlatform() {
	case appusermgrconst.ThirdGithub:
		githubAuth := oauth.NewGitHubAuth(conf)
		thirdUser, err := githubAuth.GetUserInfo(in.Code)
		if err != nil {
			return &npool.AuthLoginResponse{}, err
		}
		thirdUser.AppID = in.GetAppID()
		thirdUser.ThirdId = conf.ClientID
		thirdUser.Third = appusermgrconst.ThirdGithub
		return Login(ctx, thirdUser)
	case appusermgrconst.ThirdGoogle:
		googleAuth := oauth.NewGoogleAuth(conf)
		thirdUser, err := googleAuth.GetUserInfo(in.Code)
		if err != nil {
			return &npool.AuthLoginResponse{}, err
		}
		thirdUser.AppID = in.GetAppID()
		thirdUser.ThirdId = conf.ClientID
		thirdUser.Third = appusermgrconst.ThirdGoogle
		return Login(ctx, thirdUser)
	}
	return &npool.AuthLoginResponse{}, nil
}

func Login(ctx context.Context, thirdUser *appusermgrpb.AppUserThird) (*npool.AuthLoginResponse, error) {
	tUser, err := grpc2.GetAppUserThirdByAppThird(ctx, &appusermgrpb.GetAppUserThirdByAppThirdRequest{
		AppID:       thirdUser.AppID,
		ThirdID:     thirdUser.ThirdId,
		ThirdUserID: thirdUser.ThirdUserId,
	})
	if err != nil {
		return &npool.AuthLoginResponse{}, err
	}

	userID := ""
	if tUser != nil {
		userID = tUser.UserID
	} else {
		user, err := grpc2.CreateAppUserWithThird(ctx, &appusermgrpb.CreateAppUserWithThirdRequest{
			User: &appusermgrpb.AppUser{
				AppID: thirdUser.AppID,
			},
			Third: thirdUser,
		})
		if err != nil {
			return &npool.AuthLoginResponse{}, err
		}
		userID = user.ID
	}

	meta, err := MetadataFromContext(ctx)
	if err != nil {
		return &npool.AuthLoginResponse{}, fmt.Errorf("fail create login metadata: %v", err)
	}
	meta.AppID, err = uuid.Parse(thirdUser.AppID)
	if err != nil {
		return &npool.AuthLoginResponse{}, err
	}
	userInfo, err := grpc2.GetAppUserInfo(ctx, &appusermgrpb.GetAppUserInfoRequest{
		ID: userID,
	})
	if err != nil {
		return &npool.AuthLoginResponse{}, err
	}
	meta.UserID = uuid.MustParse(userInfo.User.ID)
	meta.ThirdUserID = thirdUser.ThirdUserId
	meta.Third = thirdUser.Third
	token, err := createToken(meta)
	if err != nil {
		return &npool.AuthLoginResponse{}, fmt.Errorf("fail create token: %v", err)
	}
	return &npool.AuthLoginResponse{
		Info:  userInfo,
		Token: token,
	}, err
}
