package tx_flood

import (
	"time"

	"github.com/urfave/cli"

	"github.com/evrynet-official/evrynet-tools/accounts"
	"github.com/evrynet-official/evrynet-tools/lib/node"
)

const (
	numTxPerAccFlag                = "num-tx-per-acc"
	floodModeFlag                  = "flood-mode"
	continuousFlooding             = "continuous"
	sleepDurationBetweenFloodsFlag = "sleep-duration"
)

// NewTxFloodFlags return flags to tx flood
func NewTxFloodFlags() []cli.Flag {
	flags := []cli.Flag{
		cli.IntFlag{
			Name:  numTxPerAccFlag,
			Usage: "Number of transactions want to use for an account",
			Value: 1,
		}, cli.IntFlag{
			Name:  floodModeFlag,
			Usage: "Flood mode when send Tx: 0: Random, 1: Normal Tx, 2: Tx with SC",
			Value: 0,
		},
		cli.BoolFlag{
			Name:  continuousFlooding,
			Usage: "Flood continuously if set to true",
		},
		cli.DurationFlag{
			Name:  sleepDurationBetweenFloodsFlag,
			Usage: "Time to sleep after each batch of numAccount*numTxPerAcc flooding",
			Value: time.Second,
		},
	}
	flags = append(flags, node.NewEvrynetNodeFlags()...)
	return flags
}

// NewTxFloodFromFlags will send tx flood
func NewTxFloodFromFlags(ctx *cli.Context) (tf *TxFlood, err error) {
	tf = &TxFlood{
		NumAcc:        ctx.Int(accounts.NumAccountsFlag.Name),
		NumTxPerAcc:   ctx.Int(numTxPerAccFlag),
		Seed:          ctx.String(accounts.SeedFlag.Name),
		FloodMode:     FloodMode(ctx.Int(floodModeFlag)),
		Continuous:    ctx.Bool(continuousFlooding),
		SleepInterval: ctx.Duration(sleepDurationBetweenFloodsFlag),
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
