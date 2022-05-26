package auth

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/thirdlogingateway"
	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db/ent"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db/ent/auth"
	"github.com/google/uuid"
)

type Auth struct {
	*db.Entity
}

func New(ctx context.Context, tx *ent.Tx) (*Auth, error) {
	e, err := db.NewEntity(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("fail create entity: %v", err)
	}

	return &Auth{
		Entity: e,
	}, nil
}

func (s *Auth) rowToObject(row *ent.Auth) *npool.Auth {
	return &npool.Auth{
		ID:           row.ID.String(),
		AppID:        row.AppID.String(),
		ThirdPartyID: row.ThirdPartyID.String(),
		AppKey:       row.AppKey,
		AppSecret:    row.AppSecret,
		RedirectURL:  row.RedirectURL,
	}
}

func (s *Auth) Create(ctx context.Context, in *npool.Auth) (*npool.Auth, error) {
	var info *ent.Auth
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Auth.Create().
			SetAppID(uuid.MustParse(in.GetAppID())).
			SetThirdPartyID(uuid.MustParse(in.GetThirdPartyID())).
			SetAppKey(in.GetAppKey()).
			SetAppSecret(in.GetAppSecret()).
			SetRedirectURL(in.GetRedirectURL()).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create auth: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Auth) CreateBulk(ctx context.Context, in []*npool.Auth) ([]*npool.Auth, error) {
	rows := []*ent.Auth{}
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		bulk := make([]*ent.AuthCreate, len(in))
		for i, info := range in {
			bulk[i] = s.Tx.Auth.Create().
				SetAppID(uuid.MustParse(info.GetAppID())).
				SetThirdPartyID(uuid.MustParse(info.GetThirdPartyID())).
				SetAppKey(info.GetAppKey()).
				SetAppSecret(info.GetAppSecret()).
				SetRedirectURL(info.GetRedirectURL())
		}
		rows, err = s.Tx.Auth.CreateBulk(bulk...).Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create bulk auth: %v", err)
	}

	infos := []*npool.Auth{}
	for _, row := range rows {
		infos = append(infos, s.rowToObject(row))
	}

	return infos, nil
}

func (s *Auth) Update(ctx context.Context, in *npool.Auth) (*npool.Auth, error) {
	var info *ent.Auth
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Auth.UpdateOneID(uuid.MustParse(in.GetID())).
			SetAppID(uuid.MustParse(in.GetAppID())).
			SetThirdPartyID(uuid.MustParse(in.GetThirdPartyID())).
			SetAppKey(in.GetAppKey()).
			SetAppSecret(in.GetAppSecret()).
			SetRedirectURL(in.GetRedirectURL()).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail update auth: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Auth) Rows(ctx context.Context, conds cruder.Conds, offset, limit int) ([]*npool.Auth, int, error) {
	rows := []*ent.Auth{}
	var total int

	err := db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		total, err = stm.Count(_ctx)
		if err != nil {
			return fmt.Errorf("fail count auth: %v", err)
		}

		rows, err = stm.
			Offset(offset).
			Order(ent.Desc(auth.FieldUpdatedAt)).
			Limit(limit).
			All(_ctx)
		if err != nil {
			return fmt.Errorf("fail query auth: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail get auth: %v", err)
	}

	infos := []*npool.Auth{}
	for _, row := range rows {
		infos = append(infos, s.rowToObject(row))
	}

	return infos, total, nil
}

func (s *Auth) RowOnly(ctx context.Context, conds cruder.Conds) (*npool.Auth, error) {
	var info *ent.Auth

	err := db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		info, err = stm.Only(_ctx)
		if err != nil {
			return fmt.Errorf("fail query auth: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get auth: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Auth) queryFromConds(conds cruder.Conds) (*ent.AuthQuery, error) {
	stm := s.Tx.Auth.Query()
	for k, v := range conds {
		switch k {
		case constant.FieldID:
			id, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid id: %v", err)
			}
			stm = stm.Where(auth.ID(id))
		case constant.AuthFieldAppID:
			id, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid app id: %v", err)
			}
			stm = stm.Where(auth.ID(id))
		default:
			return nil, fmt.Errorf("invalid auth field")
		}
	}
	return stm, nil
}
