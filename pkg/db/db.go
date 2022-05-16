package db

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/go-service-framework/pkg/mysql"
	"github.com/NpoolPlatform/third-login-gateway/pkg/db/ent"

	// runtime
	_ "github.com/NpoolPlatform/third-login-gateway/pkg/db/ent/runtime"
)

func client() (*ent.Client, error) {
	conn, err := mysql.GetConn()
	if err != nil {
		return nil, err
	}
	drv := entsql.OpenDB(dialect.MySQL, conn)
	return ent.NewClient(ent.Driver(drv)), nil
}

func Init() error {
	cli, err := client()
	if err != nil {
		return err
	}
	return cli.Schema.Create(context.Background())
}

func Client() (*ent.Client, error) {
	return client()
}

func WithTx(ctx context.Context, tx *ent.Tx, fn func(ctx context.Context) error) error {
	succ := false
	defer func() {
		if !succ {
			err := tx.Rollback()
			if err != nil {
				logger.Sugar().Errorf("fail rollback: %v", err)
				return
			}
		}
	}()
	if err := fn(ctx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %v", err)
	}
	succ = true
	return nil
}

func Do(ctx context.Context, fn func(ctx context.Context, cli *ent.Client) error) error {
	var timeOut = 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	cli, err := Client()
	if err != nil {
		return fmt.Errorf("fail get db client: %v", err)
	}

	if err := fn(ctx, cli); err != nil {
		return err
	}
	return nil
}

type Entity struct {
	Tx *ent.Tx
}

func NewEntity(ctx context.Context, _tx *ent.Tx) (*Entity, error) {
	if _tx != nil {
		return &Entity{
			Tx: _tx,
		}, nil
	}

	cli, err := Client()
	if err != nil {
		return nil, fmt.Errorf("fail get db client: %v", err)
	}
	_tx, err = cli.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail get client transaction: %v", err)
	}

	return &Entity{
		Tx: _tx,
	}, nil
}
