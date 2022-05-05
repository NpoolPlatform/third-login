package notification

import (
	"context"
	"time"

	npool "github.com/NpoolPlatform/message/npool/notification"
	"github.com/NpoolPlatform/notification/pkg/db"
	"github.com/NpoolPlatform/notification/pkg/db/ent"
	"github.com/NpoolPlatform/notification/pkg/db/ent/template"

	"github.com/google/uuid"

	"golang.org/x/xerrors"
)

const (
	dbTimeout = 5 * time.Second
)

func validateTemplate(info *npool.Template) error {
	if _, err := uuid.Parse(info.GetAppID()); err != nil {
		return xerrors.Errorf("invalid app id: %v", err)
	}
	if _, err := uuid.Parse(info.GetLangID()); err != nil {
		return xerrors.Errorf("invalid lang id: %v", err)
	}
	if info.GetTitle() == "" {
		return xerrors.Errorf("invalid title")
	}
	if info.GetContent() == "" {
		return xerrors.Errorf("invalid content")
	}
	return nil
}

func dbRowToTemplate(row *ent.Template) *npool.Template {
	return &npool.Template{
		ID:       row.ID.String(),
		AppID:    row.AppID.String(),
		LangID:   row.LangID.String(),
		UsedFor:  row.UsedFor,
		Title:    row.Title,
		Content:  row.Content,
		CreateAt: row.CreateAt,
	}
}

func Create(ctx context.Context, in *npool.CreateTemplateRequest) (*npool.CreateTemplateResponse, error) {
	if err := validateTemplate(in.GetInfo()); err != nil {
		return nil, xerrors.Errorf("invalid parameter: %v", err)
	}

	cli, err := db.Client()
	if err != nil {
		return nil, xerrors.Errorf("fail get db client: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	info, err := cli.
		Template.
		Create().
		SetAppID(uuid.MustParse(in.GetInfo().GetAppID())).
		SetLangID(uuid.MustParse(in.GetInfo().GetLangID())).
		SetUsedFor(in.GetInfo().GetUsedFor()).
		SetTitle(in.GetInfo().GetTitle()).
		SetContent(in.GetInfo().GetContent()).
		Save(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail create notification: %v", err)
	}

	return &npool.CreateTemplateResponse{
		Info: dbRowToTemplate(info),
	}, nil
}

func Get(ctx context.Context, in *npool.GetTemplateRequest) (*npool.GetTemplateResponse, error) {
	id, err := uuid.Parse(in.GetID())
	if err != nil {
		return nil, xerrors.Errorf("invalid template id: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	cli, err := db.Client()
	if err != nil {
		return nil, xerrors.Errorf("fail get db client: %v", err)
	}

	infos, err := cli.
		Template.
		Query().
		Where(
			template.ID(id),
		).
		All(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail query app sms template: %v", err)
	}

	var templates *npool.Template

	for _, info := range infos {
		templates = dbRowToTemplate(info)
		break
	}

	return &npool.GetTemplateResponse{
		Info: templates,
	}, nil
}

func Update(ctx context.Context, in *npool.UpdateTemplateRequest) (*npool.UpdateTemplateResponse, error) {
	id, err := uuid.Parse(in.GetInfo().GetID())
	if err != nil {
		return nil, xerrors.Errorf("invalid template id: %v", err)
	}

	if err := validateTemplate(in.GetInfo()); err != nil {
		return nil, xerrors.Errorf("invalid parameter: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	cli, err := db.Client()
	if err != nil {
		return nil, xerrors.Errorf("fail get db client: %v", err)
	}

	info, err := cli.
		Template.
		UpdateOneID(id).
		SetTitle(in.GetInfo().GetTitle()).
		SetContent(in.GetInfo().GetContent()).
		Save(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail update template: %v", err)
	}

	return &npool.UpdateTemplateResponse{
		Info: dbRowToTemplate(info),
	}, nil
}

func GetByApp(ctx context.Context, in *npool.GetTemplatesByAppRequest) (*npool.GetTemplatesByAppResponse, error) {
	appID, err := uuid.Parse(in.GetAppID())
	if err != nil {
		return nil, xerrors.Errorf("invalid app id: %v", err)
	}
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	cli, err := db.Client()
	if err != nil {
		return nil, xerrors.Errorf("fail get db client: %v", err)
	}

	infos, err := cli.
		Template.
		Query().
		Where(
			template.AppID(appID),
		).
		All(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail query app sms template: %v", err)
	}

	templates := []*npool.Template{}
	for _, info := range infos {
		templates = append(templates, dbRowToTemplate(info))
	}

	return &npool.GetTemplatesByAppResponse{
		Infos: templates,
	}, nil
}

func GetByAppLangUsedFor(ctx context.Context, in *npool.GetTemplateByAppLangUsedForRequest) (*npool.GetTemplateByAppLangUsedForResponse, error) {
	appID, err := uuid.Parse(in.GetAppID())
	if err != nil {
		return nil, xerrors.Errorf("invalid app id: %v", err)
	}

	langID, err := uuid.Parse(in.GetLangID())
	if err != nil {
		return nil, xerrors.Errorf("invalid lang id: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	cli, err := db.Client()
	if err != nil {
		return nil, xerrors.Errorf("fail get db client: %v", err)
	}

	infos, err := cli.
		Template.
		Query().
		Where(
			template.And(
				template.AppID(appID),
				template.LangID(langID),
				template.UsedFor(in.GetUsedFor()),
			),
		).
		All(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fail query app sms template: %v", err)
	}

	var templates *npool.Template
	for _, info := range infos {
		templates = dbRowToTemplate(info)
		break
	}

	return &npool.GetTemplateByAppLangUsedForResponse{
		Info: templates,
	}, nil
}
