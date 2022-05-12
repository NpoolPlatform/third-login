package platform

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/third-login-gateway"
	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db/ent"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db/ent/platform"
	"github.com/google/uuid"
)

type Platform struct {
	*db.Entity
}

func New(ctx context.Context, tx *ent.Tx) (*Platform, error) {
	e, err := db.NewEntity(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("fail create entity: %v", err)
	}

	return &Platform{
		Entity: e,
	}, nil
}

func (s *Platform) rowToObject(row *ent.Platform) *npool.Platform {
	return &npool.Platform{
		ID:                row.ID.String(),
		AppID:             row.AppID.String(),
		Platform:          row.Platform,
		PlatformAppKey:    row.PlatformAppKey,
		PlatformAppSecret: row.PlatformAppSecret,
		LogoUrl:           row.LogoURL,
		RedirectUrl:       row.RedirectURL,
	}
}

func (s *Platform) Create(ctx context.Context, in *npool.Platform) (*npool.Platform, error) {
	var info *ent.Platform
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Platform.Create().
			SetAppID(uuid.MustParse(in.GetAppID())).
			SetPlatform(in.GetPlatform()).
			SetPlatformAppKey(in.GetPlatformAppKey()).
			SetPlatformAppSecret(in.GetPlatformAppSecret()).
			SetLogoURL(in.GetLogoUrl()).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create stock: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Platform) Update(ctx context.Context, in *npool.Platform) (*npool.Platform, error) {
	var info *ent.Platform
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Platform.UpdateOneID(uuid.MustParse(in.GetID())).
			SetAppID(uuid.MustParse(in.GetAppID())).
			SetPlatform(in.GetPlatform()).
			SetPlatformAppKey(in.GetPlatformAppKey()).
			SetPlatformAppSecret(in.GetPlatformAppSecret()).
			SetLogoURL(in.GetLogoUrl()).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail update stock: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Platform) Rows(ctx context.Context, conds cruder.Conds, offset, limit int) ([]*npool.Platform, int, error) {
	rows := []*ent.Platform{}
	var total int

	err := db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		total, err = stm.Count(_ctx)
		if err != nil {
			return fmt.Errorf("fail count platform: %v", err)
		}

		rows, err = stm.
			Offset(offset).
			Order(ent.Desc("updated_at")).
			Limit(limit).
			All(_ctx)
		if err != nil {
			return fmt.Errorf("fail query platform: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail get platform: %v", err)
	}

	infos := []*npool.Platform{}
	for _, row := range rows {
		infos = append(infos, s.rowToObject(row))
	}

	return infos, total, nil
}

func (s *Platform) queryFromConds(conds cruder.Conds) (*ent.PlatformQuery, error) { //nolint
	stm := s.Tx.Platform.Query()
	for k, v := range conds {
		switch k {
		case constant.FieldID:
			id, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid id: %v", err)
			}
			stm = stm.Where(platform.ID(id))
		case constant.PlatformFieldAppID:
			val, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid value type: %v", err)
			}
			stm = stm.Where(platform.AppID(val))
		case constant.PlatformFieldPlatform:
			val, err := cruder.AnyTypeString(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid value type: %v", err)
			}
			stm = stm.Where(platform.Platform(val))
		default:
			return nil, fmt.Errorf("invalid platform field")
		}
	}

	return stm, nil
}
