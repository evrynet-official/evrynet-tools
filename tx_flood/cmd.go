package tx_flood

import (
	"time"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/accounts"
	"github.com/evrynet-official/evrynet-tools/lib/node"
)

const (
	numTxPerAccFlag     = "num-tx-per-acc"
	floodModeFlag       = "flood-mode"
	nonstopFlag         = "nonstop"
	nonstopDurationFlag = "nonstop-duration"
)

// NewTxFloodFlags return flags to tx flood
func NewTxFloodFlags() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			Name:  numTxPerAccFlag,
			Usage: "Number of transactions want to use for an account",
			Value: 1,
		}, cli.IntFlag{
			Name:  floodModeFlag,
			Usage: "Flood mode when send Tx: 0: Random, 1: Normal Tx, 2: Tx with SC",
			Value: int(DefaultMode),
		}, cli.BoolFlag{
			Name:  nonstopFlag,
			Usage: "To enable flood tx continuously",
		}, cli.DurationFlag{
			Name:  nonstopDurationFlag,
			Usage: "Time to start new flood",
			Value: time.Second,
		}}
}

// NewTxFloodFromFlags will send tx flood
func NewTxFloodFromFlags(ctx *cli.Context) (tf *TxFlood, err error) {
	tf = &TxFlood{
		NumAcc:          ctx.Int(accounts.NumAccountsFlag.Name),
		NumTxPerAcc:     ctx.Int(numTxPerAccFlag),
		Seed:            ctx.String(accounts.SeedFlag.Name),
		FloodMode:       FloodMode(ctx.Int(floodModeFlag)),
		Nonstop:         ctx.Bool(nonstopFlag),
		NonstopDuration: ctx.Duration(nonstopDurationFlag),
	}

	tf.Accounts, err = accounts.GenerateAccounts(tf.NumAcc, tf.Seed)
	if err != nil {
		return nil, err
	}

	tf.EvrClient, err = node.NewEvrynetClientFromFlags(ctx)
	if err != nil {
		return nil, err
	}
	return tf, nil
}
