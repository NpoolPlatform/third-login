package readuser

import (
	"context"
	"time"

	npool "github.com/NpoolPlatform/message/npool/notification"
	"github.com/NpoolPlatform/notification/pkg/db"
	"github.com/NpoolPlatform/notification/pkg/db/ent"
	"github.com/NpoolPlatform/notification/pkg/db/ent/readuser"

	"github.com/google/uuid"

	"golang.org/x/xerrors"
)

const (
	dbTimeout = 5 * time.Second
)

func validateReadUser(info *npool.ReadUser) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		return xerrors.Errorf("invalid app id: %v", err)
	}
	if _, err := uuid.Parse(info.GetUserID()); err != nil {
		return xerrors.Errorf("invalid user id: %v", err)
	}
	if _, err := uuid.Parse(info.GetAnnouncementID()); err != nil {
		return xerrors.Errorf("invalid announcement id: %v", err)
	}
	return nil
}

func dbRowToReadUser(row *ent.ReadUser) *npool.ReadUser {
	return &npool.ReadUser{
		ID:             row.ID.String(),
		AppID:          row.AppID.String(),
		UserID:         row.UserID.String(),
		AnnouncementID: row.AnnouncementID.String(),
		CreateAt:       row.CreateAt,
	}
}

func Create(ctx context.Context, in *npool.CreateReadUserRequest) (*npool.CreateReadUserResponse, error) {
	if err := validateReadUser(in.GetInfo()); err != nil {
		return nil, xerrors.Errorf("invalid parameter: %v", err)
	}

	cli, err := db.Client()
	if err != nil {
		return nil, xerrors.Errorf("fail get db client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	info, err := cli.
		ReadUser.
		Create().
		SetAppID(uuid.MustParse(in.GetInfo().GetAppID())).
		SetUserID(uuid.MustParse(in.GetInfo().GetUserID())).
		SetAnnouncementID(uuid.MustParse(in.GetInfo().GetAnnouncementID())).
		Save(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail create readuser: %v", err)
	}

	return &npool.CreateReadUserResponse{
		Info: dbRowToReadUser(info),
	}, nil
}

func Check(ctx context.Context, in *npool.CheckReadUserRequest) (*npool.CheckReadUserResponse, error) {
	if err := validateReadUser(in.GetInfo()); err != nil {
		return nil, xerrors.Errorf("invlaid parameter: %v", err)
	}

	cli, err := db.Client()
	if err != nil {
		return nil, xerrors.Errorf("fail get db client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	infos, err := cli.
		ReadUser.
		Query().
		Where(
			readuser.And(
				readuser.AppID(uuid.MustParse(in.GetInfo().GetAppID())),
				readuser.UserID(uuid.MustParse(in.GetInfo().GetUserID())),
				readuser.AnnouncementID(uuid.MustParse(in.GetInfo().GetAnnouncementID())),
			),
		).
		All(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail query readuser: %v", err)
	}

	var myReadUser *npool.ReadUser
	for _, info := range infos {
		myReadUser = dbRowToReadUser(info)
		break
	}

	return &npool.CheckReadUserResponse{
		Info: myReadUser,
	}, nil
}
