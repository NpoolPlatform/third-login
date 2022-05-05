package mailbox

import (
	"context"
	"time"

	npool "github.com/NpoolPlatform/message/npool/notification"
	"github.com/NpoolPlatform/notification/pkg/db"
	"github.com/NpoolPlatform/notification/pkg/db/ent"

	"github.com/google/uuid"

	"golang.org/x/xerrors"
)

const (
	dbTimeout = 5 * time.Second
)

func validateMail(info *npool.Mail) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		return xerrors.Errorf("invalid app id: %v", err)
	}
	if _, err := uuid.Parse(info.GetFromUserID()); err != nil {
		return xerrors.Errorf("invalid from user id: %v", err)
	}
	if _, err := uuid.Parse(info.GetToUserID()); err != nil {
		return xerrors.Errorf("invalid to user id: %v", err)
	}
	if info.GetTitle() == "" {
		return xerrors.Errorf("invalid title")
	}
	if info.GetContent() == "" {
		return xerrors.Errorf("invalid content")
	}
	return nil
}

func dbRowToMail(row *ent.MailBox) *npool.Mail {
	return &npool.Mail{
		ID:         row.ID.String(),
		AppID:      row.AppID.String(),
		FromUserID: row.FromUserID.String(),
		ToUserID:   row.ToUserID.String(),
		Title:      row.Title,
		Content:    row.Content,
		CreateAt:   row.CreateAt,
	}
}

func CreateMail(ctx context.Context, in *npool.CreateMailRequest) (*npool.CreateMailResponse, error) {
	if err := validateMail(in.GetInfo()); err != nil {
		return nil, xerrors.Errorf("invalid parameter: %v", err)
	}

	cli, err := db.Client()
	if err != nil {
		return nil, xerrors.Errorf("fail get db client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	info, err := cli.
		MailBox.
		Create().
		SetAppID(uuid.MustParse(in.GetInfo().GetAppID())).
		SetFromUserID(uuid.MustParse(in.GetInfo().GetFromUserID())).
		SetToUserID(uuid.MustParse(in.GetInfo().GetToUserID())).
		SetTitle(in.GetInfo().GetTitle()).
		SetContent(in.GetInfo().GetContent()).
		SetAlreadyRead(false).
		Save(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail create mail: %v", err)
	}

	return &npool.CreateMailResponse{
		Info: dbRowToMail(info),
	}, nil
}

func UpdateMail(ctx context.Context, in *npool.UpdateMailRequest) (*npool.UpdateMailResponse, error) {
	if err := validateMail(in.GetInfo()); err != nil {
		return nil, xerrors.Errorf("invlaid parameter: %v", err)
	}

	id, err := uuid.Parse(in.GetInfo().GetID())
	if err != nil {
		return nil, xerrors.Errorf("invalid id: %v", err)
	}

	cli, err := db.Client()
	if err != nil {
		return nil, xerrors.Errorf("fail get db client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	info, err := cli.
		MailBox.
		UpdateOneID(id).
		SetTitle(in.GetInfo().GetTitle()).
		SetContent(in.GetInfo().GetContent()).
		SetAlreadyRead(in.GetInfo().GetAlreadyRead()).
		Save(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail query mailbox: %v", err)
	}

	return &npool.UpdateMailResponse{
		Info: dbRowToMail(info),
	}, nil
}
