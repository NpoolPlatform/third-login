package notification

import (
	"context"
	"time"

	npool "github.com/NpoolPlatform/message/npool/notification"
	"github.com/NpoolPlatform/notification/pkg/db"
	"github.com/NpoolPlatform/notification/pkg/db/ent"
	"github.com/NpoolPlatform/notification/pkg/db/ent/notification"

	"github.com/google/uuid"

	"golang.org/x/xerrors"
)

const (
	dbTimeout = 5 * time.Second
)

func validateNotification(info *npool.UserNotification) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		return xerrors.Errorf("invalid app id: %v", err)
	}
	if _, err := uuid.Parse(info.GetUserID()); err != nil {
		return xerrors.Errorf("invalid user id: %v", err)
	}
	if info.GetTitle() == "" {
		return xerrors.Errorf("invalid title")
	}
	if info.GetContent() == "" {
		return xerrors.Errorf("invalid content")
	}
	return nil
}

func dbRowToNotification(row *ent.Notification) *npool.UserNotification {
	return &npool.UserNotification{
		ID:       row.ID.String(),
		AppID:    row.AppID.String(),
		UserID:   row.UserID.String(),
		Title:    row.Title,
		Content:  row.Content,
		CreateAt: row.CreateAt,
	}
}

func Create(ctx context.Context, in *npool.CreateNotificationRequest) (*npool.CreateNotificationResponse, error) {
	if err := validateNotification(in.GetInfo()); err != nil {
		return nil, xerrors.Errorf("invalid parameter: %v", err)
	}

	cli, err := db.Client()
	if err != nil {
		return nil, xerrors.Errorf("fail get db client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	info, err := cli.
		Notification.
		Create().
		SetAppID(uuid.MustParse(in.GetInfo().GetAppID())).
		SetUserID(uuid.MustParse(in.GetInfo().GetUserID())).
		SetTitle(in.GetInfo().GetTitle()).
		SetContent(in.GetInfo().GetContent()).
		SetAlreadyRead(false).
		Save(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail create notification: %v", err)
	}

	return &npool.CreateNotificationResponse{
		Info: dbRowToNotification(info),
	}, nil
}

func Update(ctx context.Context, in *npool.UpdateNotificationRequest) (*npool.UpdateNotificationResponse, error) {
	id, err := uuid.Parse(in.GetInfo().GetID())
	if err != nil {
		return nil, xerrors.Errorf("invalid id: %v", err)
	}

	if err = validateNotification(in.GetInfo()); err != nil {
		return nil, xerrors.Errorf("invlaid parameter: %v", err)
	}

	cli, err := db.Client()
	if err != nil {
		return nil, xerrors.Errorf("fail get db client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	info, err := cli.
		Notification.
		UpdateOneID(id).
		SetTitle(in.GetInfo().GetTitle()).
		SetContent(in.GetInfo().GetContent()).
		SetAlreadyRead(in.GetInfo().GetAlreadyRead()).
		Save(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail update notification: %v", err)
	}

	return &npool.UpdateNotificationResponse{
		Info: dbRowToNotification(info),
	}, nil
}

func GetNotificationsByAppUser(ctx context.Context, in *npool.GetNotificationsByAppUserRequest) (*npool.GetNotificationsByAppUserResponse, error) {
	cli, err := db.Client()
	if err != nil {
		return nil, xerrors.Errorf("fail get db client: %v", err)
	}

	appID, err := uuid.Parse(in.GetAppID())
	if err != nil {
		return nil, xerrors.Errorf("invalid app id: %v", err)
	}

	userID, err := uuid.Parse(in.GetUserID())
	if err != nil {
		return nil, xerrors.Errorf("invalid user id: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	infos, err := cli.
		Notification.
		Query().
		Where(
			notification.And(
				notification.AppID(appID),
				notification.UserID(userID),
			),
		).
		All(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail query notification: %v", err)
	}

	notifications := []*npool.UserNotification{}
	for _, info := range infos {
		notifications = append(notifications, dbRowToNotification(info))
	}

	return &npool.GetNotificationsByAppUserResponse{
		Infos: notifications,
	}, nil
}
