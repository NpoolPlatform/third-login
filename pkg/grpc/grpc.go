package grpc

import (
	"context"
	"fmt"
	"time"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/message/const"
	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
)

const (
	grpcTimeout = 10 * time.Second
)

type handle func(_ctx context.Context, cli appusermgrpb.AppUserManagerClient) (cruder.Any, error)

func doAppUser(ctx context.Context, fn handle) (cruder.Any, error) {
	_ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	return fn(_ctx, cli)
}

func GetAppUserThirdByAppThird(ctx context.Context,
	in *appusermgrpb.GetAppUserThirdPartyByAppThirdPartyIDRequest) (*appusermgrpb.AppUserThirdParty,
	error) {
	info, err := doAppUser(ctx, func(_ctx context.Context, cli appusermgrpb.AppUserManagerClient) (cruder.Any, error) {
		resp, err := cli.GetAppUserThirdPartyByAppThirdPartyID(ctx, in)
		if err != nil {
			return nil, fmt.Errorf("fail get app user third: %v", err)
		}
		return resp.Info, nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get app user third: %v", err)
	}
	return info.(*appusermgrpb.AppUserThirdParty), nil
}

func CreateAppUserWithThird(ctx context.Context, in *appusermgrpb.CreateAppUserWithThirdPartyRequest) (*appusermgrpb.AppUser, error) {
	info, err := doAppUser(ctx, func(_ctx context.Context, cli appusermgrpb.AppUserManagerClient) (cruder.Any, error) {
		resp, err := cli.CreateAppUserWithThirdParty(ctx, in)
		if err != nil {
			return nil, fmt.Errorf("fail create app user with third: %v", err)
		}
		return resp.Info, nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail create app user with third: %v", err)
	}
	return info.(*appusermgrpb.AppUser), nil
}

func GetAppUserInfo(ctx context.Context, in *appusermgrpb.GetAppUserInfoRequest) (*appusermgrpb.AppUserInfo, error) {
	info, err := doAppUser(ctx, func(_ctx context.Context, cli appusermgrpb.AppUserManagerClient) (cruder.Any, error) {
		resp, err := cli.GetAppUserInfo(ctx, in)
		if err != nil {
			return nil, fmt.Errorf("fail get app user info: %v", err)
		}
		return resp.Info, nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get app user info: %v", err)
	}
	return info.(*appusermgrpb.AppUserInfo), nil
}
