package transactions

import "github.com/urfave/cli"

var (
	numAccountsFlag = cli.IntFlag{
		Name:  "numacc",
		Usage: "Number of accounts want to use for flood tx",
		Value: 4,
	}
	numTxPerAccFlag = cli.IntFlag{
		Name:  "numtxperacc",
		Usage: "Number of transactions want to use for an account",
		Value: 1,
	}
	seedFlag = cli.StringFlag{
		Name:  "seed",
		Usage: "Seed to generate private key account",
		Value: "evrynet",
	}
	rpcEndpointFlag = cli.StringFlag{
		Name:  "rpcendpoint",
		Usage: "RPC endpoint to send request",
		Value: "http://0.0.0.0:22001",
	}
)

// NewTxFloodFlags return flags to tx flood
func NewTxFloodFlags() []cli.Flag {
	return []cli.Flag{numAccountsFlag, numTxPerAccFlag, seedFlag, rpcEndpointFlag}
}

// SendTxFlood will send tx flood
func SendTxFlood(ctx *cli.Context) error {
	txFlood := TxFlood{
		NumAcc:      ctx.Int(numAccountsFlag.Name),
		NumTxPerAcc: ctx.Int(numTxPerAccFlag.Name),
		Seed:        ctx.String(seedFlag.Name),
		RPCEndpoint: ctx.String(rpcEndpointFlag.Name),
	}
	if err := txFlood.floodTx(); err != nil {
		return err
	}
	return nil
}
