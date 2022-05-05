package notification

import (
	"context"
	"strings"

	npool "github.com/NpoolPlatform/message/npool/notification"
	constant "github.com/NpoolPlatform/notification/pkg/const"
	crud "github.com/NpoolPlatform/notification/pkg/crud/notification"
	templatecrud "github.com/NpoolPlatform/notification/pkg/crud/template"
	"golang.org/x/xerrors"
)

func CreateNotification(ctx context.Context, in *npool.CreateNotificationRequest) (*npool.CreateNotificationResponse, error) {
	template, err := templatecrud.GetByAppLangUsedFor(ctx, &npool.GetTemplateByAppLangUsedForRequest{
		AppID:   in.GetInfo().GetAppID(),
		LangID:  in.GetLangID(),
		UsedFor: in.GetUsedFor(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get template: %v", err)
	}
	if template.GetInfo() == nil {
		return nil, xerrors.Errorf("fail get template")
	}
	if in.GetMessage() != "" {
		template.Info.Content = strings.ReplaceAll(template.Info.Content, constant.MessageTemplate, in.GetMessage())
	}
	if in.GetUserName() != "" {
		template.Info.Content = strings.ReplaceAll(template.Info.Content, constant.NameTemplate, in.GetUserName())
	}

	resp, err := crud.Create(ctx, &npool.CreateNotificationRequest{
		Info: &npool.UserNotification{
			AppID:   in.GetInfo().GetAppID(),
			UserID:  in.GetInfo().GetUserID(),
			Title:   template.GetInfo().GetTitle(),
			Content: template.GetInfo().GetContent(),
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create notification: %v", err)
	}
	return resp, nil
}
