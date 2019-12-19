package tx_flood

import (
	"github.com/evrynet-official/evrynet-client/ethclient"
	"github.com/evrynet-official/evrynet-tools/accounts"
	"github.com/urfave/cli"
)

const (
	rpcEndpointFlag    = "rpcendpoint"
	defaultRPCEndpoint = "http://0.0.0.0:22001"
	numTxPerAccFlag    = "num-tx-per-acc"
	floodModeFlag      = "flood-mode"
)

// NewTxFloodFlags return flags to tx flood
func NewTxFloodFlags() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			Name:  numTxPerAccFlag,
			Usage: "Number of transactions want to use for an account",
			Value: 1,
		}, cli.StringFlag{
			Name:  "rpcendpoint",
			Usage: "RPC endpoint to send request",
			Value: defaultRPCEndpoint,
		}, cli.IntFlag{
			Name:  floodModeFlag,
			Usage: "Flood mode when send Tx: 0: Random, 1: Normal Tx, 2: Tx with SC",
			Value: 0,
		}}
}

// NewEthereumClientFromFlag returns Ethereum client from flag variable, or error if occurs
func NewEthereumClientFromFlag(ctx *cli.Context) (*ethclient.Client, error) {
	ethereumNodeURL := ctx.String(rpcEndpointFlag)
	return ethclient.Dial(ethereumNodeURL)
}

// NewTxFloodFromFlags will send tx flood
func NewTxFloodFromFlags(ctx *cli.Context) (tf *TxFlood, err error) {
	tf = &TxFlood{
		NumAcc:      ctx.Int(accounts.NumAccountsFlag.Name),
		NumTxPerAcc: ctx.Int(numTxPerAccFlag),
		Seed:        ctx.String(accounts.SeedFlag.Name),
		FloodMode:   FloodMode(ctx.Int(floodModeFlag)),
	}

	tf.Accounts, err = accounts.GenerateAccounts(tf.NumAcc, tf.Seed)
	if err != nil {
		return nil, err
	}

	tf.EvrClient, err = NewEthereumClientFromFlag(ctx)
	if err != nil {
		return nil, err
	}
	return tf, nil
}
