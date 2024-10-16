package mid

import (
	"context"
	"errors"

	"github.com/PavelDonchenko/ecommerce-micro/common/database/sqldb"
)

type ctxKey int

const (
	trKey ctxKey = iota + 1
)

func setTran(ctx context.Context, tx sqldb.CommitRollbacker) context.Context {
	return context.WithValue(ctx, trKey, tx)
}

// GetTran retrieves the value that can manage a transaction.
func GetTran(ctx context.Context) (sqldb.CommitRollbacker, error) {
	v, ok := ctx.Value(trKey).(sqldb.CommitRollbacker)
	if !ok {
		return nil, errors.New("transaction not found in context")
	}

	return v, nil
}
