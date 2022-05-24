package thirdauth

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/thirdlogingateway"
	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db/ent"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db/ent/thirdauth"
	"github.com/google/uuid"
)

type ThirdAuth struct {
	*db.Entity
}

func New(ctx context.Context, tx *ent.Tx) (*ThirdAuth, error) {
	e, err := db.NewEntity(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("fail create entity: %v", err)
	}

	return &ThirdAuth{
		Entity: e,
	}, nil
}

func (s *ThirdAuth) rowToObject(row *ent.ThirdAuth) *npool.ThirdAuth {
	return &npool.ThirdAuth{
		ID:             row.ID.String(),
		AppID:          row.AppID.String(),
		Third:          row.Third,
		ThirdAppKey:    row.ThirdAppKey,
		ThirdAppSecret: row.ThirdAppSecret,
		LogoUrl:        row.LogoURL,
		RedirectUrl:    row.RedirectURL,
	}
}

func (s *ThirdAuth) Create(ctx context.Context, in *npool.ThirdAuth) (*npool.ThirdAuth, error) {
	var info *ent.ThirdAuth
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.ThirdAuth.Create().
			SetAppID(uuid.MustParse(in.GetAppID())).
			SetThird(in.GetThird()).
			SetThirdAppKey(in.GetThirdAppKey()).
			SetThirdAppSecret(in.GetThirdAppSecret()).
			SetLogoURL(in.GetLogoUrl()).
			SetRedirectURL(in.GetRedirectUrl()).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create third auth: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *ThirdAuth) CreateBulk(ctx context.Context, in []*npool.ThirdAuth) ([]*npool.ThirdAuth, error) {
	rows := []*ent.ThirdAuth{}
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		bulk := make([]*ent.ThirdAuthCreate, len(in))
		for i, info := range in {
			bulk[i] = s.Tx.ThirdAuth.Create().
				SetAppID(uuid.MustParse(info.GetAppID())).
				SetThird(info.GetThird()).
				SetThirdAppKey(info.GetThirdAppKey()).
				SetThirdAppSecret(info.GetThirdAppSecret()).
				SetLogoURL(info.GetLogoUrl()).
				SetRedirectURL(info.GetRedirectUrl())
		}
		rows, err = s.Tx.ThirdAuth.CreateBulk(bulk...).Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create third auth: %v", err)
	}

	infos := []*npool.ThirdAuth{}
	for _, row := range rows {
		infos = append(infos, s.rowToObject(row))
	}

	return infos, nil
}

func (s *ThirdAuth) Update(ctx context.Context, in *npool.ThirdAuth) (*npool.ThirdAuth, error) {
	var info *ent.ThirdAuth
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.ThirdAuth.UpdateOneID(uuid.MustParse(in.GetID())).
			SetAppID(uuid.MustParse(in.GetAppID())).
			SetThird(in.GetThird()).
			SetThirdAppKey(in.GetThirdAppKey()).
			SetThirdAppSecret(in.GetThirdAppSecret()).
			SetLogoURL(in.GetLogoUrl()).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail update third auth: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *ThirdAuth) Rows(ctx context.Context, conds cruder.Conds, offset, limit int) ([]*npool.ThirdAuth, int, error) {
	rows := []*ent.ThirdAuth{}
	var total int

	err := db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		total, err = stm.Count(_ctx)
		if err != nil {
			return fmt.Errorf("fail count third auth: %v", err)
		}

		rows, err = stm.
			Offset(offset).
			Order(ent.Desc("updated_at")).
			Limit(limit).
			All(_ctx)
		if err != nil {
			return fmt.Errorf("fail query third auth: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail get third auth: %v", err)
	}

	infos := []*npool.ThirdAuth{}
	for _, row := range rows {
		infos = append(infos, s.rowToObject(row))
	}

	return infos, total, nil
}

func (s *ThirdAuth) RowOnly(ctx context.Context, conds cruder.Conds) (*npool.ThirdAuth, error) {
	var info *ent.ThirdAuth

	err := db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		info, err = stm.Only(_ctx)
		if err != nil {
			return fmt.Errorf("fail query third auth: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get third auth: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *ThirdAuth) queryFromConds(conds cruder.Conds) (*ent.ThirdAuthQuery, error) {
	stm := s.Tx.ThirdAuth.Query()
	for k, v := range conds {
		switch k {
		case constant.FieldID:
			id, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid id: %v", err)
			}
			stm = stm.Where(thirdauth.ID(id))
		case constant.ThirdAuthFieldAppID:
			val, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid value app id: %v", err)
			}
			stm = stm.Where(thirdauth.AppID(val))
		case constant.ThirdAuthFieldThird:
			val, err := cruder.AnyTypeString(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid value third: %v", err)
			}
			stm = stm.Where(thirdauth.Third(val))
		default:
			return nil, fmt.Errorf("invalid third auth field")
		}
	}

	return stm, nil
}
