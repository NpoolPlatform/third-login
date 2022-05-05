package readuser

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

func assertReadUser(t *testing.T, actual, expected *npool.ReadUser) {
	assert.Equal(t, actual.AppID, expected.AppID)
	assert.Equal(t, actual.UserID, expected.UserID)
	assert.Equal(t, actual.AnnouncementID, expected.AnnouncementID)
	assert.NotEqual(t, actual.CreateAt, 0)
}

func TestCRUD(t *testing.T) {
	readuser := npool.ReadUser{
		AppID:          uuid.New().String(),
		UserID:         uuid.New().String(),
		AnnouncementID: uuid.New().String(),
	}

	resp, err := Create(context.Background(), &npool.CreateReadUserRequest{
		Info: &readuser,
	})
	if assert.Nil(t, err) {
		assert.NotEqual(t, resp.Info.ID, uuid.UUID{}.String())
		assertReadUser(t, resp.Info, &readuser)
	}

	resp1, err := Check(context.Background(), &npool.CheckReadUserRequest{
		Info: &readuser,
	})
	if assert.Nil(t, err) {
		assert.Equal(t, resp1.Info.ID, resp.Info.ID)
		assertReadUser(t, resp1.Info, &readuser)
	}
}
