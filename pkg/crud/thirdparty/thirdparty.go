package thirdparty

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/thirdlogingateway"
	constant "github.com/NpoolPlatform/third-login-gateway/pkg/const"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db/ent"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db/ent/thirdparty"
	"github.com/google/uuid"
)

type ThirdParty struct {
	*db.Entity
}

func New(ctx context.Context, tx *ent.Tx) (*ThirdParty, error) {
	e, err := db.NewEntity(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("fail create entity: %v", err)
	}

	return &ThirdParty{
		Entity: e,
	}, nil
}

func (s *ThirdParty) rowToObject(row *ent.ThirdParty) *npool.ThirdParty {
	return &npool.ThirdParty{
		ID:        row.ID.String(),
		BrandName: row.BrandName,
		Logo:      row.Logo,
		Domain:    row.Domain,
	}
}

func (s *ThirdParty) Create(ctx context.Context, in *npool.ThirdParty) (*npool.ThirdParty, error) {
	var info *ent.ThirdParty
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.ThirdParty.Create().
			SetBrandName(in.GetBrandName()).
			SetLogo(in.GetLogo()).
			SetDomain(in.GetDomain()).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create third party: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *ThirdParty) CreateBulk(ctx context.Context, in []*npool.ThirdParty) ([]*npool.ThirdParty, error) {
	rows := []*ent.ThirdParty{}
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		bulk := make([]*ent.ThirdPartyCreate, len(in))
		for i, info := range in {
			bulk[i] = s.Tx.ThirdParty.Create().
				SetBrandName(info.GetBrandName()).
				SetLogo(info.GetLogo()).
				SetDomain(info.GetDomain())
		}
		rows, err = s.Tx.ThirdParty.CreateBulk(bulk...).Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create bulk third party: %v", err)
	}

	infos := []*npool.ThirdParty{}
	for _, row := range rows {
		infos = append(infos, s.rowToObject(row))
	}

	return infos, nil
}

func (s *ThirdParty) Update(ctx context.Context, in *npool.ThirdParty) (*npool.ThirdParty, error) {
	var info *ent.ThirdParty
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.ThirdParty.UpdateOneID(uuid.MustParse(in.GetID())).
			SetBrandName(in.GetBrandName()).
			SetLogo(in.GetLogo()).
			SetDomain(in.GetDomain()).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail update third party: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *ThirdParty) Rows(ctx context.Context, conds cruder.Conds, offset, limit int) ([]*npool.ThirdParty, int, error) {
	rows := []*ent.ThirdParty{}
	var total int

	err := db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		total, err = stm.Count(_ctx)
		if err != nil {
			return fmt.Errorf("fail count third party: %v", err)
		}

		rows, err = stm.
			Offset(offset).
			Order(ent.Desc(thirdparty.FieldUpdatedAt)).
			Limit(limit).
			All(_ctx)
		if err != nil {
			return fmt.Errorf("fail query third party: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail get third party: %v", err)
	}

	infos := []*npool.ThirdParty{}
	for _, row := range rows {
		infos = append(infos, s.rowToObject(row))
	}

	return infos, total, nil
}

func (s *ThirdParty) Row(ctx context.Context, id uuid.UUID) (*npool.ThirdParty, error) {
	var info *ent.ThirdParty
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.ThirdParty.Query().Where(thirdparty.ID(id)).Only(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail get third party: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *ThirdParty) RowOnly(ctx context.Context, conds cruder.Conds) (*npool.ThirdParty, error) {
	var info *ent.ThirdParty

	err := db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		info, err = stm.Only(_ctx)
		if err != nil {
			return fmt.Errorf("fail query third party: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get third party: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *ThirdParty) queryFromConds(conds cruder.Conds) (*ent.ThirdPartyQuery, error) {
	stm := s.Tx.ThirdParty.Query()
	for k, v := range conds {
		switch k {
		case constant.FieldID:
			id, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid id: %v", err)
			}
			stm = stm.Where(thirdparty.ID(id))
		default:
			return nil, fmt.Errorf("invalid third party field")
		}
	}
	return stm, nil
}
