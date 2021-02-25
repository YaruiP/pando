package sys

import (
	"context"

	"github.com/fox-one/pando/core"
	"github.com/fox-one/pando/pkg/maker"
	"github.com/fox-one/pando/pkg/uuid"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func HandleWithdraw(wallets core.WalletStore) maker.HandlerFunc {
	return func(ctx context.Context, r *maker.Request) error {
		if err := require(r.Gov, "not-authorized"); err != nil {
			return err
		}

		var (
			asset    uuid.UUID
			amount   decimal.Decimal
			opponent uuid.UUID
		)

		if err := require(r.Scan(&asset, &amount, &opponent) == nil, "bad-data"); err != nil {
			return err
		}

		memo := core.TransferAction{
			ID:     r.TraceID(),
			Source: r.Action.String(),
		}.Encode()

		t := &core.Transfer{
			CreatedAt: r.Now(),
			TraceID:   uuid.Modify(r.TraceID(), memo),
			AssetID:   asset.String(),
			Amount:    decimal.Decimal{},
			Memo:      memo,
			Threshold: 1,
			Opponents: []string{opponent.String()},
		}

		if err := wallets.CreateTransfers(ctx, []*core.Transfer{t}); err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("wallets.CreateTransfers")
			return err
		}

		return nil
	}
}
