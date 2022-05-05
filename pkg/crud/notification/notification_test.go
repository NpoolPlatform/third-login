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

func assertUserNotification(t *testing.T, actual, expected *npool.UserNotification) {
	assert.Equal(t, actual.AppID, expected.AppID)
	assert.Equal(t, actual.UserID, expected.UserID)
	assert.Equal(t, actual.Title, expected.Title)
	assert.Equal(t, actual.Content, expected.Content)
	assert.NotEqual(t, actual.CreateAt, 0)
}

func TestCRUD(t *testing.T) {
	notification := npool.UserNotification{
		AppID:   uuid.New().String(),
		UserID:  uuid.New().String(),
		Title:   "Haaaaaaaaaaaaaaaaaaaa",
		Content: "Coooooooooooooooooooooooooooooooo",
	}

	resp, err := Create(context.Background(), &npool.CreateNotificationRequest{
		Info: &notification,
	})
	if assert.Nil(t, err) {
		assert.NotEqual(t, resp.Info.ID, uuid.UUID{}.String())
		assertUserNotification(t, resp.Info, &notification)
	}

	notification.ID = resp.Info.ID
	notification.Title = notification.Content + "hhhhhhhh"
	notification.Content = notification.Title + "Ccccccccc"

	resp1, err := Update(context.Background(), &npool.UpdateNotificationRequest{
		Info: &notification,
	})
	if assert.Nil(t, err) {
		assert.Equal(t, resp1.Info.ID, notification.ID)
		assertUserNotification(t, resp1.Info, &notification)
	}

	resp2, err := GetNotificationsByAppUser(context.Background(), &npool.GetNotificationsByAppUserRequest{
		AppID:  notification.AppID,
		UserID: notification.UserID,
	})
	if assert.Nil(t, err) {
		assert.NotEqual(t, len(resp2.Infos), 0)
	}
}
