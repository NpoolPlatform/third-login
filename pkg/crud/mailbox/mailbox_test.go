package mailbox

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

func assertMail(t *testing.T, actual, expected *npool.Mail) {
	assert.Equal(t, actual.AppID, expected.AppID)
	assert.Equal(t, actual.FromUserID, expected.FromUserID)
	assert.Equal(t, actual.ToUserID, expected.ToUserID)
	assert.Equal(t, actual.Title, expected.Title)
	assert.Equal(t, actual.Content, expected.Content)
	assert.NotEqual(t, actual.CreateAt, 0)
}

func TestCRUD(t *testing.T) {
	mail := npool.Mail{
		AppID:      uuid.New().String(),
		FromUserID: uuid.New().String(),
		ToUserID:   uuid.New().String(),
		Title:      "Haaaaaaaaaaaaaaaaaaaa",
		Content:    "Coooooooooooooooooooooooooooooooo",
	}

	resp, err := CreateMail(context.Background(), &npool.CreateMailRequest{
		Info: &mail,
	})
	if assert.Nil(t, err) {
		assert.NotEqual(t, resp.Info.ID, uuid.UUID{}.String())
		assertMail(t, resp.Info, &mail)
	}

	mail.ID = resp.Info.ID
	mail.Title = mail.Content + "hhhhhhhh"
	mail.Content = mail.Title + "Ccccccccc"

	resp1, err := UpdateMail(context.Background(), &npool.UpdateMailRequest{
		Info: &mail,
	})
	if assert.Nil(t, err) {
		assert.Equal(t, resp1.Info.ID, mail.ID)
		assertMail(t, resp1.Info, &mail)
	}
}
