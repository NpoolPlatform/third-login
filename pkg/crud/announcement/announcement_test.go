package announcement

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

func assertAnnouncement(t *testing.T, actual, expected *npool.Announcement) {
	assert.Equal(t, actual.AppID, expected.AppID)
	assert.Equal(t, actual.Title, expected.Title)
	assert.Equal(t, actual.Content, expected.Content)
	assert.NotEqual(t, actual.CreateAt, 0)
}

func TestCRUD(t *testing.T) {
	announcement := npool.Announcement{
		AppID:   uuid.New().String(),
		Title:   "Haaaaaaaaaaaaaaaaaaaa",
		Content: "Coooooooooooooooooooooooooooooooo",
	}

	resp, err := Create(context.Background(), &npool.CreateAnnouncementRequest{
		Info: &announcement,
	})
	if assert.Nil(t, err) {
		assert.NotEqual(t, resp.Info.ID, uuid.UUID{}.String())
		assertAnnouncement(t, resp.Info, &announcement)
	}

	announcement.ID = resp.Info.ID
	announcement.Title = announcement.Content + "hhhhhhhh"
	announcement.Content = announcement.Title + "Ccccccccc"

	resp1, err := Update(context.Background(), &npool.UpdateAnnouncementRequest{
		Info: &announcement,
	})
	if assert.Nil(t, err) {
		assert.Equal(t, resp1.Info.ID, announcement.ID)
		assertAnnouncement(t, resp1.Info, &announcement)
	}

	resp2, err := GetAnnouncementsByApp(context.Background(), &npool.GetAnnouncementsByAppRequest{
		AppID: announcement.AppID,
	})
	if assert.Nil(t, err) {
		assert.NotEqual(t, len(resp2.Infos), 0)
	}
}
