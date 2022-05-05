package notification

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	npool "github.com/NpoolPlatform/message/npool/notification"
	"github.com/NpoolPlatform/notification/pkg/test-init" //nolint

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func init() {
	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); err == nil && runByGithubAction {
		return
	}
	if err := testinit.Init(); err != nil {
		fmt.Printf("cannot init test stub: %v\n", err)
	}
}

func assertTemplate(t *testing.T, actual, expected *npool.Template) {
	assert.Equal(t, actual.AppID, expected.AppID)
	assert.Equal(t, actual.LangID, expected.LangID)
	assert.Equal(t, actual.UsedFor, expected.UsedFor)
	assert.Equal(t, actual.Title, expected.Title)
	assert.Equal(t, actual.Content, expected.Content)
	assert.NotEqual(t, actual.CreateAt, 0)
}

func TestCRUD(t *testing.T) {
	template := npool.Template{
		AppID:   uuid.New().String(),
		LangID:  uuid.New().String(),
		UsedFor: "test used for",
		Title:   "test title",
		Content: "test Content",
	}

	resp, err := Create(context.Background(), &npool.CreateTemplateRequest{
		Info: &template,
	})
	if assert.Nil(t, err) {
		assert.NotEqual(t, resp.Info.ID, uuid.UUID{}.String())
		assertTemplate(t, resp.Info, &template)
	}

	resp1, err := Get(context.Background(), &npool.GetTemplateRequest{
		ID: resp.Info.ID,
	})
	if assert.Nil(t, err) {
		assert.Equal(t, resp1.Info.ID, resp.Info.ID)
		assertTemplate(t, resp1.Info, &template)
	}

	template.ID = resp.Info.ID
	template.Title = template.Content + "update"
	template.Content = template.Title + "update"

	resp2, err := Update(context.Background(), &npool.UpdateTemplateRequest{
		Info: &template,
	})
	if assert.Nil(t, err) {
		assert.Equal(t, resp2.Info.ID, resp.Info.ID)
		assertTemplate(t, resp2.Info, &template)
	}

	template.ID = resp.Info.ID
	resp3, err := GetByAppLangUsedFor(context.Background(), &npool.GetTemplateByAppLangUsedForRequest{
		AppID:   template.AppID,
		LangID:  template.LangID,
		UsedFor: template.UsedFor,
	})
	if assert.Nil(t, err) {
		assert.Equal(t, resp3.Info.ID, resp.Info.ID)
		assertTemplate(t, resp3.Info, &template)
	}
}
