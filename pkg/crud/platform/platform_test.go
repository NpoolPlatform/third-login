package platform

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/third-login-gateway"
	"github.com/NpoolPlatform/third-login-gateway/pkg/test-init" //nolint

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
	platform := npool.Platform{
		AppID:             uuid.New().String(),
		Platform:          "github",
		PlatformAppKey:    "AAAAAAAAAaAaaaaaaaaaaaAAAAAAAAAAAAAAAA",
		PlatformAppSecret: "AAAAAAAAAaAaaaaaaaaaaaAAAAAAAAAAAAAAAA",
		LogoUrl:           "AAAAAAAAAAAAAAAAAAAAAAAAAAa",
	}

	schema, err := New(context.Background(), nil)
	assert.Nil(t, err)

	info, err := schema.Create(context.Background(), &platform)
	if assert.Nil(t, err) {
		if assert.NotEqual(t, info.ID, uuid.UUID{}.String()) {
			platform.ID = info.ID
		}
		if assert.NotEqual(t, info.AppID, uuid.UUID{}.String()) {
			platform.AppID = info.AppID
		}
		assert.Equal(t, info, &platform)
	}

	platform.ID = info.ID
	platform.AppID = info.AppID

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.Update(context.Background(), &platform)
	if assert.Nil(t, err) {
		assert.Equal(t, info, &platform)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	//info, err = schema.Row(context.Background(), uuid.MustParse(info.ID))
	//if assert.Nil(t, err) {
	//	assert.Equal(t, info, &stock)
	//}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	infos, total, err := schema.Rows(context.Background(),
		cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID),
		0, 0)
	if assert.Nil(t, err) {
		assert.Equal(t, total, 1)
		assert.Equal(t, infos[0], &platform)
	}

	//schema, err = New(context.Background(), nil)
	//assert.Nil(t, err)
	//
	//info, err = schema.RowOnly(context.Background(),
	//	cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID))
	//if assert.Nil(t, err) {
	//	assert.Equal(t, info, &stock)
	//}
	//
	//schema, err = New(context.Background(), nil)
	//assert.Nil(t, err)
	//
	//exist, err := schema.ExistConds(context.Background(),
	//	cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID),
	//)
	//if assert.Nil(t, err) {
	//	assert.Equal(t, exist, true)
	//}
	//
	//schema, err = New(context.Background(), nil)
	//assert.Nil(t, err)
	//
	//stock.Total = 10000
	//
	//info, err = schema.UpdateFields(context.Background(),
	//	uuid.MustParse(info.ID),
	//	cruder.NewFields().
	//		WithField(constant.StockFieldTotal, stock.Total),
	//)
	//if assert.Nil(t, err) {
	//	assert.Equal(t, info, &stock)
	//}
	//
	//schema, err = New(context.Background(), nil)
	//assert.Nil(t, err)
	//
	//stock.InService = 2
	//stock.Locked = 1
	//stock.Sold = 2
	//
	//info, err = schema.AddFields(context.Background(),
	//	uuid.MustParse(info.ID),
	//	cruder.NewFields().
	//		WithField(constant.StockFieldInService, 2).
	//		WithField(constant.StockFieldLocked, 1),
	//)
	//if assert.Nil(t, err) {
	//	assert.Equal(t, info, &stock)
	//}
	//
	//assert.Nil(t, err)
	//assert.NotNil(t, info)
	//
	//schema, err = New(context.Background(), nil)
	//assert.Nil(t, err)
	//
	//stock.InService = 1
	//stock.Locked = 0
	//
	//info, err = schema.AddFields(context.Background(),
	//	uuid.MustParse(info.ID),
	//	cruder.NewFields().
	//		WithField(constant.StockFieldInService, -1).
	//		WithField(constant.StockFieldLocked, -1),
	//)
	//if assert.Nil(t, err) {
	//	assert.Equal(t, info, &stock)
	//}
	//
	//schema, err = New(context.Background(), nil)
	//assert.Nil(t, err)
	//
	//_, err = schema.AddFields(context.Background(),
	//	uuid.MustParse(info.ID),
	//	cruder.NewFields().
	//		WithField(constant.StockFieldInService, 5002).
	//		WithField(constant.StockFieldLocked, 5001),
	//)
	//assert.NotNil(t, err)
	//
	//stock1 := &npool.Stock{
	//	GoodID:    uuid.New().String(),
	//	Total:     1000,
	//	InService: 0,
	//	Sold:      0,
	//}
	//stock2 := &npool.Stock{
	//	GoodID:    uuid.New().String(),
	//	Total:     1000,
	//	InService: 0,
	//	Sold:      0,
	//}
	//
	//schema, err = New(context.Background(), nil)
	//assert.Nil(t, err)
	//
	//infos, err = schema.CreateBulk(context.Background(), []*npool.Stock{stock1, stock2})
	//if assert.Nil(t, err) {
	//	assert.Equal(t, len(infos), 2)
	//	assert.NotEqual(t, infos[0].ID, uuid.UUID{}.String())
	//	assert.NotEqual(t, infos[1].ID, uuid.UUID{}.String())
	//}
	//
	//schema, err = New(context.Background(), nil)
	//assert.Nil(t, err)
	//
	//count, err := schema.Count(context.Background(),
	//	cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID),
	//)
	//if assert.Nil(t, err) {
	//	assert.Equal(t, count, uint32(1))
	//}
	//
	//schema, err = New(context.Background(), nil)
	//assert.Nil(t, err)
	//
	//info, err = schema.Delete(context.Background(), uuid.MustParse(info.ID))
	//if assert.Nil(t, err) {
	//	assert.Equal(t, info, &stock)
	//}
	//
	//schema, err = New(context.Background(), nil)
	//assert.Nil(t, err)
	//
	//count, err = schema.Count(context.Background(),
	//	cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID),
	//)
	//if assert.Nil(t, err) {
	//	assert.Equal(t, count, uint32(0))
	//}
	//
	//schema, err = New(context.Background(), nil)
	//assert.Nil(t, err)
	//
	//_, err = schema.Row(context.Background(), uuid.MustParse(info.ID))
	//assert.NotNil(t, err)
}
