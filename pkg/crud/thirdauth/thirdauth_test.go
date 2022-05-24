package thirdauth

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/thirdlogingateway"
	testinit "github.com/NpoolPlatform/third-login-gateway/pkg/test-init"

	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"

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

func TestCRUD(t *testing.T) {
	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); err == nil && runByGithubAction {
		return
	}

	thirdAuth := npool.ThirdAuth{
		AppID:          uuid.New().String(),
		Third:          uuid.New().String(),
		ThirdAppKey:    uuid.New().String(),
		ThirdAppSecret: uuid.New().String(),
		LogoUrl:        uuid.New().String(),
		RedirectUrl:    uuid.New().String(),
	}

	schema, err := New(context.Background(), nil)
	assert.Nil(t, err)

	info, err := schema.Create(context.Background(), &thirdAuth)
	if assert.Nil(t, err) {
		if assert.NotEqual(t, info.ID, uuid.UUID{}.String()) {
			thirdAuth.ID = info.ID
		}
		assert.Equal(t, info, &thirdAuth)
	}

	thirdAuth.ID = info.ID

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.Update(context.Background(), &thirdAuth)
	if assert.Nil(t, err) {
		assert.Equal(t, info, &thirdAuth)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	infos, total, err := schema.Rows(context.Background(),
		cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID),
		0, 0)
	if assert.Nil(t, err) {
		assert.Equal(t, total, 1)
		assert.Equal(t, infos[0], &thirdAuth)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.RowOnly(context.Background(),
		cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID))
	if assert.Nil(t, err) {
		assert.Equal(t, info, &thirdAuth)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	thirdAuth1 := &npool.ThirdAuth{
		AppID:          uuid.New().String(),
		Third:          uuid.New().String(),
		ThirdAppKey:    uuid.New().String(),
		ThirdAppSecret: uuid.New().String(),
		LogoUrl:        uuid.New().String(),
		RedirectUrl:    uuid.New().String(),
	}

	thirdAuth2 := &npool.ThirdAuth{
		AppID:          uuid.New().String(),
		Third:          uuid.New().String(),
		ThirdAppKey:    uuid.New().String(),
		ThirdAppSecret: uuid.New().String(),
		LogoUrl:        uuid.New().String(),
		RedirectUrl:    uuid.New().String(),
	}

	infos, err = schema.CreateBulk(context.Background(), []*npool.ThirdAuth{thirdAuth1, thirdAuth2})
	if assert.Nil(t, err) {
		assert.Equal(t, len(infos), 2)
		assert.NotEqual(t, infos[0].ID, uuid.UUID{}.String())
		assert.NotEqual(t, infos[1].ID, uuid.UUID{}.String())
	}
}
