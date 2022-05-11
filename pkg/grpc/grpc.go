package grpc

import (
	"context"
	"fmt"
	"time"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/message/const" //nolint
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
)

const (
	grpcTimeout = 60 * time.Second
)

func GetAppUserThirdByAppThird(ctx context.Context, in *appusermgrpb.GetAppUserThirdByAppThirdRequest) (*appusermgrpb.AppUserThird, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppUserThirdByAppThird(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail update cache: %v", err)
	}

	return resp.Info, nil
}

func CreateAppUserWithThird(ctx context.Context, in *appusermgrpb.CreateAppUserWithThirdRequest) (*appusermgrpb.AppUser, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateAppUserWithThird(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail update cache: %v", err)
	}

	return resp.Info, nil
}

func GetAppUserInfo(ctx context.Context, in *appusermgrpb.GetAppUserInfoRequest) (*appusermgrpb.AppUserInfo, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppUserInfo(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail update cache: %v", err)
	}

	return resp.Info, nil
}
